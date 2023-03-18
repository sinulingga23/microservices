# Build docker image
docker build -f Dockerfile -t <image-tag>

# Run docker image as container
docker run -p 8081:8081 -e 8081 -d --name <container-name> -it <image-tag>