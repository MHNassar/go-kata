package main

import (
	"context"
	"time"
)

type ProfileService interface {
	GetProfile(ctx context.Context, userID int) (string, error)
}

type OrderService interface {
	GetOrders(ctx context.Context, userID int) (string, error)
}

type MockProfileService struct {
	Delay time.Duration
	Err   error
}

func (m *MockProfileService) GetProfile(ctx context.Context, userID int) (string, error) {
	select {
	case <-time.After(m.Delay):
		if m.Err != nil {
			return "", m.Err
		}
		return "Name: Alice", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

type MockOrderService struct {
	Delay time.Duration
	Err   error
}

func (m *MockOrderService) GetOrders(ctx context.Context, userID int) (string, error) {
	select {
	case <-time.After(m.Delay):
		if m.Err != nil {
			return "", m.Err
		}
		return "Orders: 5", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
