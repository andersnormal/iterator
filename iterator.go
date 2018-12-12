package iterator

import (
	"fmt"
	"reflect"
)

func NewIterator(retrieve func(int, string) (string, error), bufLen func() int, retrieveBuf func() interface{}) (*Iterator, func() error) {
	i := &Iterator{
		retrieve:    retrieve,
		retrieveBuf: retrieveBuf,
		bufLen:      bufLen,
	}

	return i, i.next
}

func (i *Iterator) Len() int {
	return i.bufLen()
}

func (i *Iterator) next() error {
	i.nextCalled = true

	if i.err != nil {
		return i.err
	}

	for i.bufLen() == 0 && !i.atLast {
		if err := i.buffer(i.MaxSize); err != nil {
			i.err = err
			return i.err
		}
		if i.Token == "" {
			i.atLast = true
		}
	}

	if i.bufLen() == 0 {
		i.err = nil
	}

	return i.err
}

func (i *Iterator) buffer(size int) error {
	token, err := i.retrieve(size, i.Token)
	if err != nil {
		i.retrieveBuf()
		return err
	}

	i.Token = token
	return nil
}

func NewPager(iter Pageable, size int, token string) *Pager {
	p := &Pager{iter.Iterator(), size}
	p.iterator.Token = token

	if size < 0 {
		p.iterator.err = ErrPositiveSize
	}

	return p
}

func (p *Pager) NextPage(ptr interface{}) (string, error) {
	p.iterator.nextPageCalled = true

	if p.iterator.err != nil {
		return "", p.iterator.err
	}

	if p.iterator.bufLen() > 0 {
		return "", ErrBufferNotEmpty
	}

	bufType := reflect.PtrTo(reflect.ValueOf(p.iterator.retrieveBuf()).Type())
	if ptr == nil {
		return "", ErrNilNextPage
	}

	ptrValue := reflect.ValueOf(ptr)
	if ptrValue.Type() != bufType {
		return "", fmt.Errorf("pager: next should be of type %s, got %T", bufType, ptr)
	}

	for p.iterator.bufLen() < p.size {
		if err := p.iterator.buffer(p.size - p.iterator.bufLen()); err != nil {
			return "", p.iterator.err
		}

		if p.iterator.Token == "" {
			break
		}
	}

	e := ptrValue.Elem()
	e.Set(reflect.AppendSlice(e, reflect.ValueOf(p.iterator.retrieveBuf())))
	return p.iterator.Token, nil
}
