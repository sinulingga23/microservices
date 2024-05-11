# How to run this project.
```
docker-compose up -d
```

# How to see logs of an container.
```
docker-compose logs -f --tail=500 | grep `<container-name>`
```

# How to shutdown the services.
```
docker-compose down
```


# Common Kafka Commands
* Get list topic                  : kafka-topics --bootstrap-server=<host:port> --list
* Get topic detail                : kafka-topics --bootstrap-server=localhost:9092 --describe --topic <topic-name>
* Consumes mesages on an topic    : kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic test --from-beginning --max-messages 10

# Prometheus Metrics
## Product Service
* request_total_endpoint_add_product
* request_total_endpoint_get_products
* request_total_endpoint_get_product_by_id
* request_total_endpoint_get_products_by_ids

## Order Service
* request_total_endpoint_create_orders

# Tech Stack
* Golang
* Go-Chi
* Docker
* Kafka 
* Nginx
* Prometheus
* gRPC


