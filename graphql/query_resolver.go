package main

import (
	"context"
	"log"
	"time"

	"github.com/sunil8777/E-commerce-microservices/graphql/model"
)

type queryResolver struct {
	server *Server
}

func (r *queryResolver) Accounts(ctx context.Context, pagination *model.PaginationInput, id *string) ([]*model.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if id != nil {
		r, err := r.server.accountClient.GetAccount(ctx, *id)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		return []*model.Account{{
			ID:   r.ID,
			Name: r.Name,
		}}, nil
	}

	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = bounds(pagination)
	}

	accountList, err := r.server.accountClient.GetAccounts(ctx, take, skip)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var accounts []*model.Account
	for _, a := range accountList {
		account := &model.Account{
			ID:   a.ID,
			Name: a.Name,
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (r *queryResolver) Products(ctx context.Context, pagination *model.PaginationInput, query *string, id *string) ([]*model.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if id != nil {
		r, err := r.server.catalogClient.GetProduct(ctx, *id)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		return []*model.Product{{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Price:       r.Price,
		}}, nil
	}

	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = bounds(pagination)
	}

	q := ""
	if query != nil {
		q = *query
	}

	productList, err := r.server.catalogClient.GetProducts(ctx, nil, q, skip, take)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var products []*model.Product
	for _, a := range productList {
		products = append(products, &model.Product{
			ID:          a.ID,
			Name:        a.Name,
			Description: a.Description,
			Price:       a.Price,
		})
	}

	return products, nil
}

func bounds(p *model.PaginationInput) (uint64, uint64) {
	skipValue := uint64(0)
	takeValue := uint64(100)
	if p.Skip != nil {
		skipValue = uint64(*p.Skip)
	}
	if p.Take != nil {
		takeValue = uint64(*p.Take)
	}
	return skipValue, takeValue
}
