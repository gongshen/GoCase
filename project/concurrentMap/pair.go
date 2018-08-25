package cmap

import (
	"bytes"
	"fmt"
	"sync/atomic"
	"unsafe"
)

type linkedPair interface {
	Next() Pair
	SetNext(nextPair Pair) error
}

type Pair interface {
	linkedPair
	Key() string
	Element() interface{}
	Hash() uint64
	SetElement(element interface{}) error
	Copy() Pair
	String() string
}

type pair struct {
	key     string
	element unsafe.Pointer
	hash    uint64
	next    unsafe.Pointer
}

func newPair(key string, element interface{}) (Pair, error) {
	p := &pair{
		key:  key,
		hash: hash(key),
	}
	if element == nil {
		return nil, newIllegalParameterError("element is nil")
	}
	p.element = unsafe.Pointer(&element)
	return p, nil
}

func (p *pair) Key() string {
	return p.key
}

func (p *pair) Hash() uint64 {
	return p.hash
}

func (p *pair) Element() interface{} {
	element := atomic.LoadPointer(&p.element)
	if element == nil {
		return nil
	}
	return *(*interface{})(element)
}

func (p *pair) SetElement(element interface{}) error {
	if element == nil {
		return newIllegalParameterError("element is nil")
	}
	atomic.StorePointer(&p.element, unsafe.Pointer(&element))
	return nil
}

func (p *pair) Copy() Pair {
	pCopy, _ := newPair(p.Key(), p.Element())
	return pCopy
}

func (p *pair) String() string {
	return p.genString(false)
}

func (p *pair) Next() Pair {
	pointer := atomic.LoadPointer(&p.next)
	if pointer == nil {
		return nil
	}
	return (*pair)(pointer)
}

func (p *pair) SetNext(nextPair Pair) error {
	if nextPair == nil {
		atomic.StorePointer(&p.next, nil)
		return nil
	}
	pp, ok := nextPair.(*pair)
	if !ok {
		return newIllegalPairTypeError(nextPair)
	}
	atomic.StorePointer(&p.next, unsafe.Pointer(pp))
	return nil
}

func (p *pair) genString(detail bool) string {
	var buf bytes.Buffer
	buf.WriteString("pair{key:")
	buf.WriteString(p.Key())
	buf.WriteString(",element:")
	buf.WriteString(fmt.Sprintf("%+v", p.Element()))
	buf.WriteString(",hash:")
	buf.WriteString(fmt.Sprintf("%d", p.Hash()))
	if detail {
		buf.WriteString(",next:")
		if v := p.Next(); v != nil {
			if vv, ok := v.(*pair); ok {
				buf.WriteString(vv.genString(detail))
			} else {
				buf.WriteString("<ignore>")
			}
		}
	} else {
		buf.WriteString(",nextKey:")
		if v := p.Next(); v != nil {
			buf.WriteString(v.Key())
		}
	}
	buf.WriteString("}")
	return buf.String()
}
