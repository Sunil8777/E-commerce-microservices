package main

import "context"

type mutationResolver struct {
	server *Server
}

func (m *mutationResolver) createAccount(ctx context.Context, in AccountInput) (*Account, error) {

}

func (m *mutationResolver) createProduct(ctx context.Context, in ProductInput) (*Product, error) {

}

func (m *mutationResolver) createOrder(ctx context.Context, in OrderInput) (*Order, error) {

}
