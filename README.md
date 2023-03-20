# How to run this project.
```
docker-compose up -d
docker cp services.conf nginx:/etc/nginx/conf.d
docker exec nginx nginx -t (ensure the status is ok before go to next step)
docker exec nginx nginx -s reload
```
DONE

# How to see logs of an container.
docker-compose logs -f --tail=500 <container-name>

# How to shutdown the services.
docker-compose down
