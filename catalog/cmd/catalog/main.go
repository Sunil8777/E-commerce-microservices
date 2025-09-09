package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/sunil8777/E-commerce-microservices/catalog"
	"github.com/sunil8777/E-commerce-microservices/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r catalog.Repository

	retry.ForeverSleep(func() error {
		r, err = catalog.NewElasticRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		return err
	})

	defer r.Close()
	log.Println("listening on port 8080")
	s := catalog.NewService(r)
	log.Fatal(catalog.ListenGRPC(s, 8080))
}
