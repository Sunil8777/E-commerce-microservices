#  E-commerce Microservices

A microservices-based e-commerce platform built with **Go**, **GraphQL**, **gRPC**, **PostgreSQL**, and **Elasticsearch**.

---

##  Tech Stack

- **Language**: Go  
- **API Gateway**: GraphQL (gqlgen)  
- **Communication**: gRPC  
- **Databases**: PostgreSQL (Account, Order)  
- **Search Engine**: Elasticsearch (Catalog)  
- **Containerization**: Docker & Docker Compose  
- **ID Generation**: ksuid  

---

##  Architecture

```text
      GraphQL Gateway
             |
    -------------------------------
    |            |               |
 Account      Catalog          Order
 Client       Client           Client   (gRPC clients inside the Gateway)
    |            |               |
    v            v               v
 Account Server  Catalog Server  Order Server   (gRPC servers)
    |            |               |
 Postgres   Elasticsearch     Postgres
```
---
##  Getting Started  

1. **Clone the repository**  

```bash
git clone https://github.com/Sunil8777/E-commerce-microservices.git
cd E-commerce-microservices
```
2. **Start docker**
```bash

docker-compose up --build

