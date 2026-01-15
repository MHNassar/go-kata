package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"sync"
	"time"
)

type UserAggregator interface {
	Aggregate(ctx context.Context, id int) error
}

type Aggregator struct {
	profileSvc ProfileService
	orderSvc   OrderService
	logger     *slog.Logger
	timeout    time.Duration
}
type Option func(*Aggregator)

func WithTimeout(d time.Duration) Option {
	return func(a *Aggregator) {
		a.timeout = d
	}
}
func WithLogger(l *slog.Logger) Option {
	return func(a *Aggregator) {
		a.logger = l
	}
}
func NewAggregator(profile ProfileService, order OrderService, opts ...Option) *Aggregator {
	a := &Aggregator{
		profileSvc: profile,
		orderSvc:   order,
		logger:     slog.Default(),
		timeout:    3 * time.Second,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func (a *Aggregator) Aggregate(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()
	g, gctx := errgroup.WithContext(ctx)
	var output struct {
		result map[string]string
		mutex  sync.Mutex
	}
	output.result = make(map[string]string)
	g.Go(func() error {
		pData, err := a.profileSvc.GetProfile(gctx, id)
		if err != nil {
			a.logger.Error("failed to get profile", "error", err)
			return err
		}
		output.mutex.Lock()
		output.result["profile"] = pData
		output.mutex.Unlock()
		return nil
	})
	g.Go(func() error {
		oData, err := a.orderSvc.GetOrders(gctx, id)
		if err != nil {
			a.logger.Error("failed to get orders", "error", err)
			return err
		}
		output.mutex.Lock()
		output.result["order"] = oData
		output.mutex.Unlock()
		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("aggregation failed : %w", err)
	}
	fmt.Println(fmt.Sprintf("%s | %s", output.result["profile"], output.result["order"]))

	return nil

}

func main() {
	ctx := context.Background()
	profileSvc := &MockProfileService{Delay: 2 * time.Second}
	orderSvc := &MockOrderService{Delay: 500 * time.Millisecond}
	aggregator := NewAggregator(
		profileSvc,
		orderSvc,
		WithTimeout(3*time.Second), WithLogger(slog.Default()),
	)

	err := aggregator.Aggregate(ctx, 1)
	if err != nil {
		slog.Error("aggregation failed", "error", err)
		return
	}
	return

}
