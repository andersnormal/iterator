package iterator

import ()

// Iterator represents the state of an iterator
type Iterator struct {
	MaxSize int
	Token   string

	err         error
	atLast      bool
	retrieve    func(pageSize int, pageToken string) (nextPageToken string, err error)
	retrieveBuf func() interface{}

	bufLen func() int

	nextCalled     bool
	nextPageCalled bool
}

// Pageable is implemented by iterators that support paging.
type Pageable interface {
	Iterator() *Iterator
}

type Pager struct {
	iterator *Iterator
	size     int
}
