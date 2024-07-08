docker build . -f Dockerfile.Gorilla -t websocket_benchmark/go-gorilla && \
docker run --rm -p 8080:8080 websocket_benchmark/go-gorilla
