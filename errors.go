package iterator

import (
	"errors"
)

var (
	ErrPositiveSize   = errors.New("iterator: page size must be positive")
	ErrBufferNotEmpty = errors.New("pager: must call NextPage with an empty buffer")
	ErrNilNextPage    = errors.New("pager: nil passed to Pager.NextPage")
)
