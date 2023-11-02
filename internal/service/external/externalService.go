package external

import (
	"context"
	"time"

	"github.com/dwnGnL/testWork/internal/config"
	externalLib "github.com/dwnGnL/testWork/lib/external"
	"github.com/dwnGnL/testWork/lib/goerrors"
)

type externalClient struct {
	client            externalLib.Service
	queue             chan externalLib.Item
	writeBuffer       externalLib.Batch
	limit             uint64
	counter           int
	periodResetTicker *time.Ticker
}

func New(ctx context.Context, conf *config.Config) *externalClient {
	client := externalLib.New()
	n, p := client.GetLimits()
	var timeToReset time.Duration
	if p > time.Second {
		timeToReset = p - time.Second // reduce the time by 1 second to take advantage of the interval that has not yet been used up
	}
	c := &externalClient{
		client:            client,
		writeBuffer:       make(externalLib.Batch, 0, n),
		queue:             make(chan externalLib.Item, 1000),
		periodResetTicker: time.NewTicker(timeToReset),
		limit:             n,
	}
	c.Run(ctx)
	return c
}

func (c *externalClient) Run(ctx context.Context) {
	go c.writer(ctx)
}

func (c *externalClient) writer(ctx context.Context) {
	for {
		select {
		case object := <-c.queue:
			c.prepareForWrite(object)
			if c.canWriteNow() { //check possible of write
				c.processBatch(ctx)
			}
		case <-c.periodResetTicker.C:
			c.processBatch(ctx)
			c.counter = 0 // reset counter for other batchItem
		case <-ctx.Done():
			return
		}
	}
}

func (c *externalClient) ProcessBatch(objects externalLib.Item) {
	c.queue <- objects
}

func (c *externalClient) prepareForWrite(obj externalLib.Item) {
	c.writeBuffer = append(c.writeBuffer, obj)
}

func (c *externalClient) canWriteNow() bool {
	return len(c.writeBuffer) == int(c.limit) && c.counter == 0
}

func (c *externalClient) processBatch(ctx context.Context) {
	c.counter += len(c.writeBuffer)
	err := c.client.Process(ctx, c.writeBuffer)
	if err != nil {
		goerrors.Log().Error("process batch err, need wait before write again")
		return
	}
	c.writeBuffer = c.writeBuffer[:0]
}
