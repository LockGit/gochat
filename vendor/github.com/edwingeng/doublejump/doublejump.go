// Package doublejump provides a revamped Google's jump consistent hash.
package doublejump

import (
	"math/rand"

	"github.com/dgryski/go-jump"
)

type looseHolder struct {
	a []interface{}
	m map[interface{}]int
	f []int
}

func (this *looseHolder) add(obj interface{}) {
	if _, ok := this.m[obj]; ok {
		return
	}

	if nf := len(this.f); nf == 0 {
		this.a = append(this.a, obj)
		this.m[obj] = len(this.a) - 1
	} else {
		idx := this.f[nf-1]
		this.f = this.f[:nf-1]
		this.a[idx] = obj
		this.m[obj] = idx
	}
}

func (this *looseHolder) remove(obj interface{}) {
	if idx, ok := this.m[obj]; ok {
		this.f = append(this.f, idx)
		this.a[idx] = nil
		delete(this.m, obj)
	}
}

func (this *looseHolder) get(key uint64) interface{} {
	na := len(this.a)
	if na == 0 {
		return nil
	}

	h := jump.Hash(key, na)
	return this.a[h]
}

func (this *looseHolder) shrink() {
	if len(this.f) == 0 {
		return
	}

	var a []interface{}
	for _, obj := range this.a {
		if obj != nil {
			a = append(a, obj)
			this.m[obj] = len(a) - 1
		}
	}
	this.a = a
	this.f = nil
}

type compactHolder struct {
	a []interface{}
	m map[interface{}]int
}

func (this *compactHolder) add(obj interface{}) {
	if _, ok := this.m[obj]; ok {
		return
	}

	this.a = append(this.a, obj)
	this.m[obj] = len(this.a) - 1
}

func (this *compactHolder) shrink(a []interface{}) {
	for i, obj := range a {
		this.a[i] = obj
		this.m[obj] = i
	}
}

func (this *compactHolder) remove(obj interface{}) {
	if idx, ok := this.m[obj]; ok {
		n := len(this.a)
		this.a[idx] = this.a[n-1]
		this.m[this.a[idx]] = idx
		this.a[n-1] = nil
		this.a = this.a[:n-1]
		delete(this.m, obj)
	}
}

func (this *compactHolder) get(key uint64) interface{} {
	na := len(this.a)
	if na == 0 {
		return nil
	}

	h := jump.Hash(key*0xc6a4a7935bd1e995, na)
	return this.a[h]
}

// Hash is a revamped Google's jump consistent hash. It overcomes the shortcoming of the
// original implementation - not being able to remove nodes.
type Hash struct {
	loose   looseHolder
	compact compactHolder
}

// NewHash creates a new doublejump hash instance, which does NOT threadsafe.
func NewHash() *Hash {
	hash := &Hash{}
	hash.loose.m = make(map[interface{}]int)
	hash.compact.m = make(map[interface{}]int)
	return hash
}

// Add adds an object to the hash.
func (this *Hash) Add(obj interface{}) {
	if obj == nil {
		return
	}

	this.loose.add(obj)
	this.compact.add(obj)
}

// Remove removes an object from the hash.
func (this *Hash) Remove(obj interface{}) {
	if obj == nil {
		return
	}

	this.loose.remove(obj)
	this.compact.remove(obj)
}

// Len returns the number of objects in the hash.
func (this *Hash) Len() int {
	return len(this.compact.a)
}

// LooseLen returns the size of the inner loose object holder.
func (this *Hash) LooseLen() int {
	return len(this.loose.a)
}

// Shrink removes all empty slots from the hash.
func (this *Hash) Shrink() {
	this.loose.shrink()
	this.compact.shrink(this.loose.a)
}

// Get returns an object according to the key provided.
func (this *Hash) Get(key uint64) interface{} {
	obj := this.loose.get(key)
	switch obj {
	case nil:
		return this.compact.get(key)
	default:
		return obj
	}
}

// All returns all the objects in this Hash.
func (this *Hash) All() []interface{} {
	n := len(this.compact.a)
	if n == 0 {
		return nil
	}
	all := make([]interface{}, n)
	copy(all, this.compact.a)
	return all
}

// Random returns a random object.
func (this *Hash) Random() interface{} {
	if n := len(this.compact.a); n > 0 {
		idx := rand.Intn(n)
		return this.compact.a[idx]
	}
	return nil
}
