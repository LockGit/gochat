package etcdv3

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"go.etcd.io/etcd/clientv3"
)

const (
	// ETCDV3 backend
	ETCDV3 store.Backend = "etcdv3"
)

// EtcdV3 is the receiver type for the Store interface
type EtcdV3 struct {
	timeout time.Duration
	client  *clientv3.Client
	leaseID clientv3.LeaseID
	done    chan struct{}
}

// Register registers etcd to libkv
func Register() {
	libkv.AddStore(ETCDV3, New)
}

// New creates a new Etcd client given a list
// of endpoints and an optional tls config
func New(addrs []string, options *store.Config) (store.Store, error) {
	s := &EtcdV3{
		done: make(chan struct{}),
	}

	cfg := clientv3.Config{
		Endpoints: addrs,
	}

	if options != nil {
		s.timeout = options.ConnectionTimeout
		cfg.DialTimeout = options.ConnectionTimeout
		cfg.DialKeepAliveTimeout = options.ConnectionTimeout
		cfg.TLS = options.TLS
		cfg.Username = options.Username
		cfg.Password = options.Password

		cfg.AutoSyncInterval = 5 * time.Minute
	}
	if s.timeout == 0 {
		s.timeout = 10 * time.Second
	}

	cli, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}

	s.client = cli

	resp, err := cli.Grant(context.Background(), 30)
	if err != nil {
		cli.Close()
		return nil, err
	}
	s.leaseID = resp.ID

	go func() {
		ch, kaerr := cli.KeepAlive(context.Background(), resp.ID)
		if kaerr != nil {
			log.Fatal(kaerr)
		}
		for {
			select {
			case <-s.done:
				return
			case <-ch:
			}
		}
	}()

	return s, nil
}

func (s *EtcdV3) normalize(key string) string {
	key = store.Normalize(key)
	return strings.TrimPrefix(key, "/")
}

// Put a value at the specified key
func (s *EtcdV3) Put(key string, value []byte, options *store.WriteOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	_, err := s.client.Put(ctx, key, string(value), clientv3.WithLease(s.leaseID))
	cancel()

	return err
}

// Get a value given its key
func (s *EtcdV3) Get(key string) (*store.KVPair, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	resp, err := s.client.Get(ctx, key)
	cancel()
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, store.ErrKeyNotFound
	}

	pair := &store.KVPair{
		Key:       key,
		Value:     resp.Kvs[0].Value,
		LastIndex: uint64(resp.Kvs[0].Version),
	}

	return pair, nil
}

// Delete the value at the specified key
func (s *EtcdV3) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	_, err := s.client.Delete(ctx, key)
	cancel()

	return err
}

// Exists verifies if a Key exists in the store
func (s *EtcdV3) Exists(key string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	resp, err := s.client.Get(ctx, key)
	cancel()
	if err != nil {
		return false, err
	}

	return len(resp.Kvs) != 0, nil
}

// Watch for changes on a key
func (s *EtcdV3) Watch(key string, stopCh <-chan struct{}) (<-chan *store.KVPair, error) {
	watchCh := make(chan *store.KVPair)

	go func() {
		defer close(watchCh)

		pair, err := s.Get(key)
		if err != nil {
			return
		}
		watchCh <- pair

		rch := s.client.Watch(context.Background(), key)
		for {
			select {
			case <-s.done:
				return
			case wresp := <-rch:
				for _, event := range wresp.Events {
					watchCh <- &store.KVPair{
						Key:       string(event.Kv.Key),
						Value:     event.Kv.Value,
						LastIndex: uint64(event.Kv.Version),
					}
				}
			}
		}
	}()

	return watchCh, nil
}

// WatchTree watches for changes on child nodes under
// a given directory
func (s *EtcdV3) WatchTree(directory string, stopCh <-chan struct{}) (<-chan []*store.KVPair, error) {
	watchCh := make(chan []*store.KVPair)

	go func() {
		defer close(watchCh)

		list, err := s.List(directory)
		if err != nil {
			return
		}

		watchCh <- list

		rch := s.client.Watch(context.Background(), directory, clientv3.WithPrefix())
		for {
			select {
			case <-s.done:
				return
			case <-rch:
				list, err := s.List(directory)
				if err != nil {
					return
				}
				watchCh <- list
			}
		}
	}()

	return watchCh, nil
}

type etcdLock struct {
	session *concurrency.Session
	mutex   *concurrency.Mutex
}

// NewLock creates a lock for a given key.
// The returned Locker is not held and must be acquired
// with `.Lock`. The Value is optional.
func (s *EtcdV3) NewLock(key string, options *store.LockOptions) (store.Locker, error) {
	return nil, errors.New("not implemented")
}

// List the content of a given prefix
func (s *EtcdV3) List(directory string) ([]*store.KVPair, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	resp, err := s.client.Get(ctx, directory, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	kvpairs := make([]*store.KVPair, 0, len(resp.Kvs))

	if len(resp.Kvs) == 0 {
		return kvpairs, nil
	}

	for _, kv := range resp.Kvs {
		pair := &store.KVPair{
			Key:       string(kv.Key),
			Value:     kv.Value,
			LastIndex: uint64(kv.Version),
		}
		kvpairs = append(kvpairs, pair)
	}

	return kvpairs, nil
}

// DeleteTree deletes a range of keys under a given directory
func (s *EtcdV3) DeleteTree(directory string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	_, err := s.client.Delete(ctx, directory, clientv3.WithPrefix())
	cancel()

	return err
}

// AtomicPut CAS operation on a single value.
// Pass previous = nil to create a new key.
func (s *EtcdV3) AtomicPut(key string, value []byte, previous *store.KVPair, options *store.WriteOptions) (bool, *store.KVPair, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)

	var revision int64
	var presp *clientv3.PutResponse
	var txresp *clientv3.TxnResponse
	var err error
	if previous == nil {
		presp, err = s.client.Put(ctx, key, string(value), clientv3.WithLease(s.leaseID))
		if presp != nil {
			revision = presp.Header.GetRevision()
		}
	} else {

		var cmps = []clientv3.Cmp{
			clientv3.Compare(clientv3.Value("key"), "=", string(previous.Value)),
			clientv3.Compare(clientv3.Version("key"), "=", int64(previous.LastIndex)),
		}
		txresp, err = s.client.Txn(ctx).If(cmps...).
			Then(clientv3.OpPut(key, string(value))).
			Commit()
		if txresp != nil {
			revision = txresp.Header.GetRevision()
		}
	}
	cancel()

	if err != nil {
		return false, nil, err
	}

	pair := &store.KVPair{
		Key:       key,
		Value:     value,
		LastIndex: uint64(revision),
	}

	return true, pair, nil
}

// AtomicDelete cas deletes a single value
func (s *EtcdV3) AtomicDelete(key string, previous *store.KVPair) (bool, error) {
	var deleted = false
	var err error
	var txresp *clientv3.TxnResponse
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	if previous == nil {
		var dresp *clientv3.DeleteResponse
		dresp, err = s.client.Delete(ctx, key)
		deleted = len(dresp.PrevKvs) != 0
	} else {
		// dresp, err = s.client.Delete(ctx, key, clientv3.WithRev(int64(previous.LastIndex)))

		var cmps = []clientv3.Cmp{
			clientv3.Compare(clientv3.Value("key"), "=", string(previous.Value)),
			clientv3.Compare(clientv3.Version("key"), "=", int64(previous.LastIndex)),
		}
		txresp, err = s.client.Txn(ctx).If(cmps...).
			Then(clientv3.OpDelete(key)).
			Commit()
		deleted = txresp.Succeeded
	}
	cancel()

	if err != nil {
		return false, err
	}

	return deleted, nil
}

// Close closes the client connection
func (s *EtcdV3) Close() {
	close(s.done)
	s.client.Close()
}
