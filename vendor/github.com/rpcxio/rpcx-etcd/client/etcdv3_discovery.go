package client

import (
	"strings"
	"sync"
	"time"

	"github.com/rpcxio/libkv"
	"github.com/rpcxio/libkv/store"
	estore "github.com/rpcxio/rpcx-etcd/store"
	etcd "github.com/rpcxio/rpcx-etcd/store/etcdv3"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/log"
)

func init() {
	etcd.Register()
}

// EtcdV3Discovery is a etcd service discovery.
// It always returns the registered servers in etcd.
type EtcdV3Discovery struct {
	basePath string
	kv       store.Store
	pairsMu  sync.RWMutex
	pairs    []*client.KVPair
	chans    []chan []*client.KVPair
	mu       sync.Mutex

	// -1 means it always retry to watch until zookeeper is ok, 0 means no retry.
	RetriesAfterWatchFailed int

	filter           client.ServiceDiscoveryFilter
	AllowKeyNotFound bool

	stopCh chan struct{}
}

// NewEtcdV3Discovery returns a new EtcdV3Discovery.
func NewEtcdV3Discovery(basePath string, servicePath string, etcdAddr []string, allowKeyNotFound bool, options *store.Config) (client.ServiceDiscovery, error) {
	kv, err := libkv.NewStore(estore.ETCDV3, etcdAddr, options)
	if err != nil {
		log.Infof("cannot create store: %v", err)
		return nil, err
	}

	if ev3, ok := kv.(*etcd.EtcdV3); ok {
		ev3.AllowKeyNotFound = allowKeyNotFound
	}

	return NewEtcdV3DiscoveryStore(basePath+"/"+servicePath, kv, allowKeyNotFound)
}

// NewEtcdV3DiscoveryStore return a new EtcdV3Discovery with specified store.
func NewEtcdV3DiscoveryStore(basePath string, kv store.Store, allowKeyNotFound bool) (client.ServiceDiscovery, error) {
	if len(basePath) > 1 && strings.HasSuffix(basePath, "/") {
		basePath = basePath[:len(basePath)-1]
	}

	d := &EtcdV3Discovery{basePath: basePath, kv: kv}
	d.stopCh = make(chan struct{})

	ps, err := kv.List(basePath)
	if err != nil {
		if !allowKeyNotFound || err != store.ErrKeyNotFound {
			log.Errorf("cannot get services of from registry: %v, err: %v", basePath, err)
			return nil, err
		}
	}
	pairs := make([]*client.KVPair, 0, len(ps))
	var prefix string
	for _, p := range ps {
		if prefix == "" {
			if strings.HasPrefix(p.Key, "/") {
				if strings.HasPrefix(d.basePath, "/") {
					prefix = d.basePath + "/"
				} else {
					prefix = "/" + d.basePath + "/"
				}
			} else {
				if strings.HasPrefix(d.basePath, "/") {
					prefix = d.basePath[1:] + "/"
				} else {
					prefix = d.basePath + "/"
				}
			}
		}
		if p.Key == prefix[:len(prefix)-1] || !strings.HasPrefix(p.Key, prefix) {
			continue
		}
		k := strings.TrimPrefix(p.Key, prefix)
		pair := &client.KVPair{Key: k, Value: string(p.Value)}
		if d.filter != nil && !d.filter(pair) {
			continue
		}
		pairs = append(pairs, pair)
	}
	d.pairsMu.Lock()
	d.pairs = pairs
	d.pairsMu.Unlock()
	d.RetriesAfterWatchFailed = -1
	d.AllowKeyNotFound = allowKeyNotFound

	go d.watch()
	return d, nil
}

// NewEtcdV3DiscoveryTemplate returns a new EtcdV3Discovery template.
func NewEtcdV3DiscoveryTemplate(basePath string, etcdAddr []string, allowKeyNotFound bool, options *store.Config) (client.ServiceDiscovery, error) {
	if len(basePath) > 1 && strings.HasSuffix(basePath, "/") {
		basePath = basePath[:len(basePath)-1]
	}

	kv, err := libkv.NewStore(estore.ETCDV3, etcdAddr, options)
	if err != nil {
		log.Infof("cannot create store: %v", err)
		return nil, err
	}

	if ev3, ok := kv.(*etcd.EtcdV3); ok {
		ev3.AllowKeyNotFound = allowKeyNotFound
	}

	return NewEtcdV3DiscoveryStore(basePath, kv, allowKeyNotFound)
}

// Clone clones this ServiceDiscovery with new servicePath.
func (d *EtcdV3Discovery) Clone(servicePath string) (client.ServiceDiscovery, error) {
	// assume servicePath does not contains /
	basePath := d.basePath
	i := strings.LastIndex(basePath, "/")
	if i > 0 {
		basePath = basePath[:i]
	}
	return NewEtcdV3DiscoveryStore(basePath+"/"+servicePath, d.kv, d.AllowKeyNotFound)
}

// SetFilter sets the filer.
func (d *EtcdV3Discovery) SetFilter(filter client.ServiceDiscoveryFilter) {
	d.filter = filter
}

// GetServices returns the servers
func (d *EtcdV3Discovery) GetServices() []*client.KVPair {
	d.pairsMu.RLock()
	defer d.pairsMu.RUnlock()

	return d.pairs
}

// WatchService returns a nil chan.
func (d *EtcdV3Discovery) WatchService() chan []*client.KVPair {
	d.mu.Lock()
	defer d.mu.Unlock()

	ch := make(chan []*client.KVPair, 10)
	d.chans = append(d.chans, ch)
	return ch
}

func (d *EtcdV3Discovery) RemoveWatcher(ch chan []*client.KVPair) {
	d.mu.Lock()
	defer d.mu.Unlock()

	var chans []chan []*client.KVPair
	for _, c := range d.chans {
		if c == ch {
			continue
		}

		chans = append(chans, c)
	}

	d.chans = chans
}

func (d *EtcdV3Discovery) watch() {
	defer func() {
		d.kv.Close()
	}()

rewatch:
	for {
		var err error
		var c <-chan []*store.KVPair
		var tempDelay time.Duration

		retry := d.RetriesAfterWatchFailed
		for d.RetriesAfterWatchFailed < 0 || retry >= 0 {
			c, err = d.kv.WatchTree(d.basePath, nil)
			if err != nil {
				if d.RetriesAfterWatchFailed > 0 {
					retry--
				}
				if tempDelay == 0 {
					tempDelay = 1 * time.Second
				} else {
					tempDelay *= 2
				}
				if max := 30 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Warnf("can not watchtree (with retry %d, sleep %v): %s: %v", retry, tempDelay, d.basePath, err)
				time.Sleep(tempDelay)
				continue
			}
			break
		}

		if err != nil {
			log.Errorf("can't watch %s: %v", d.basePath, err)
			return
		}

		for {
			select {
			case <-d.stopCh:
				log.Info("discovery has been closed")
				return
			case ps, ok := <-c:
				if !ok {
					break rewatch
				}
				var pairs []*client.KVPair // latest servers
				if ps == nil && !d.AllowKeyNotFound {
					d.pairsMu.Lock()
					d.pairs = pairs
					d.pairsMu.Unlock()
					continue
				}
				var prefix string
				for _, p := range ps {
					if prefix == "" {
						if strings.HasPrefix(p.Key, "/") {
							if strings.HasPrefix(d.basePath, "/") {
								prefix = d.basePath + "/"
							} else {
								prefix = "/" + d.basePath + "/"
							}
						} else {
							if strings.HasPrefix(d.basePath, "/") {
								prefix = d.basePath[1:] + "/"
							} else {
								prefix = d.basePath + "/"
							}
						}
					}
					if p.Key == prefix[:len(prefix)-1] || !strings.HasPrefix(p.Key, prefix) {
						continue
					}

					k := strings.TrimPrefix(p.Key, prefix)
					pair := &client.KVPair{Key: k, Value: string(p.Value)}
					if d.filter != nil && !d.filter(pair) {
						continue
					}
					pairs = append(pairs, pair)
				}
				d.pairsMu.Lock()
				d.pairs = pairs
				d.pairsMu.Unlock()

				d.mu.Lock()
				for _, ch := range d.chans {
					ch := ch
					go func() {
						defer func() {
							recover()
						}()

						select {
						case ch <- pairs:
						case <-time.After(time.Minute):
							log.Warn("chan is full and new change has been dropped")
						}
					}()
				}
				d.mu.Unlock()
			}
		}

		// log.Warn("chan is closed and will rewatch")
	}
}

func (d *EtcdV3Discovery) Close() {
	close(d.stopCh)
}
