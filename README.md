# How to run this project.
```
docker-compose up -d
docker cp services.conf nginx:/etc/nginx/conf.d
docker exec nginx nginx -t (ensure the status is ok before go to the next step)
docker exec nginx nginx -s reload
```

# How to see logs of an container.
docker-compose logs -f --tail=500 <container-name>

# How to shutdown the services.
docker-compose down


# Common Kafka Commands
Get list topic                  : kafka-topics --bootstrap-server=<host:port> --list
Get topic detail                : kafka-topics --bootstrap-server=localhost:9092 --describe --topic <topic-name>
Consumes mesages on an topic    :  

kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic test --from-beginning --max-messages 10
# Good to know
On server where your admin run kafka find kafka-console-consumer.sh by command find . -name kafka-console-consumer.sh then go to that directory and run for read message from your topic
