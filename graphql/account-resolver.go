package main

import (
	"context"
	"log"
	"time"

	"github.com/sunil8777/E-commerce-microservices/graphql/model"
	"github.com/sunil8777/E-commerce-microservices/order"
)

type accountResolver struct {
	server *Server
}

func (r *accountResolver) Orders(ctx context.Context, obj *model.Account) ([]*order.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	orderList, err := r.server.orderClient.GetOrderForAccount(ctx, obj.ID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var orders []*order.Order
	for _, o := range orderList {
		var products []order.OrderedProduct
		for _, p := range o.Products {
			products = append(products, order.OrderedProduct{
				ID:          p.ID,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
				Quantity:    uint64(p.Quantity),
			})
		}
		orders = append(orders, &order.Order{
			ID:         o.ID,
			CreatedAt:  o.CreatedAt,
			TotalPrice: o.TotalPrice,
			Products:   products,
		})
	}

	return orders, nil
}
