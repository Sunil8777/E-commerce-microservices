package catalog

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
)

var (
	ErrNotFound = errors.New("entity not found")
)

type Repository interface {
	Close()
	PutProduct(ctx context.Context, p Product) error
	GetProductByID(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
}

type elasticSearchRepository struct {
	client *elasticsearch.Client
}

type productDocument struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       float64 `json:"price"`
}

func NewElasticRepository(url string) (Repository, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			url,
		},
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &elasticSearchRepository{client}, nil
}

func (r *elasticSearchRepository) Close() {
}

func (r *elasticSearchRepository) PutProduct(ctx context.Context, p Product) error {

	body, err := json.Marshal(productDocument{
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	})
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      "catalog",
		DocumentID: p.ID,
		Body:       bytes.NewReader(body),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return err
}

func (r *elasticSearchRepository) GetProductByID(ctx context.Context, id string) (*Product, error) {
	res, err := r.client.Get(
		"catalog",
		id,
		r.client.Get.WithContext(ctx),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil, ErrNotFound
	}

	var doc struct {
		Source productDocument `json:"_source"`
	}

	if err := json.NewDecoder(res.Body).Decode(&doc); err != nil {
		return nil, err
	}

	return &Product{
		ID:          id,
		Name:        doc.Source.Name,
		Description: doc.Source.Description,
		Price:       doc.Source.Price,
	}, nil
}

func (r *elasticSearchRepository) ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	query := map[string]interface{}{
		"from": skip,
		"size": take,
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}

	body, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex("catalog"),
		r.client.Search.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var esRes struct {
		Hits struct {
			Hits []struct {
				ID     string  `json:"_id"`
				Source Product `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&esRes); err != nil {
		return nil, err
	}

	products := []Product{}
	for _, hit := range esRes.Hits.Hits {
		p := hit.Source
		p.ID = hit.ID
		products = append(products, p)
	}

	return products, nil

}

func (r *elasticSearchRepository) ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error) {
	if len(ids) == 0 {
		return []Product{}, nil
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"ids": map[string]interface{}{
				"values": ids,
			},
		},
	}

	body, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex("catalog"),
		r.client.Search.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil, ErrNotFound
	}

	var sr struct {
		Hits struct {
			Hits []struct {
				ID     string  `json:"_id"`
				Source Product `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
		return nil, err
	}

	products := []Product{}
	for _, hit := range sr.Hits.Hits {
		p := hit.Source
		p.ID = hit.ID
		products = append(products, p)
	}

	return products, nil
}

func (r *elasticSearchRepository) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	if query == "" {
		return []Product{}, nil
	}

	esQuery := map[string]interface{}{
		"from": skip,
		"size": take,
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"name", "description"},
			},
		},
	}

	body, err := json.Marshal(esQuery)
	if err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex("catalog"),
		r.client.Search.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil, ErrNotFound
	}

	var sr struct {
		Hits struct {
			Hits []struct {
				ID     string  `json:"_id"`
				Source Product `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
		return nil, err
	}

	products := []Product{}
	for _, hit := range sr.Hits.Hits {
		p := hit.Source
		p.ID = hit.ID
		products = append(products, p)
	}

	return products, nil
}
