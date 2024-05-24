# Build docker image
docker build -f Dockerfile -t <image-tag>

# Run docker image as container
docker run -p 8081:8081 -e 8081 -d --name <container-name> -it <image-tag>


# Generated Protobuf
From the root service:
```
protoc --proto_path=proto --go_out=proto-generated --go-grpc_out=proto-generated --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative base.proto product.proto
```

# Commands
## Generate an key 
* openssl rand -base64 756 > keyfile_mongo.txt