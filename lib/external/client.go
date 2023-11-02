package external

import (
	"context"
	"errors"
	"log"
	"time"
)

// ErrBlocked reports if service is blocked.
var ErrBlocked = errors.New("blocked")

// Service defines external service that can process batches of items.
type Service interface {
	GetLimits() (n uint64, p time.Duration)
	Process(ctx context.Context, batch Batch) error
}

// Batch is a batch of items.
type Batch []Item

// Item is some abstract item.
type Item struct {
	count uint64
	p     time.Duration
	n     uint64
}

func New() *Item {
	return &Item{
		p: time.Minute,
		n: 10,
	}
}

func (i Item) GetLimits() (n uint64, p time.Duration) {
	return i.n, i.p
}

func (i Item) Process(ctx context.Context, batch Batch) error {
	log.Println(batch)
	if i.count > i.count+uint64(len(batch)) {
		return ErrBlocked
	}
	i.count += uint64(len(batch))

	return nil
}
