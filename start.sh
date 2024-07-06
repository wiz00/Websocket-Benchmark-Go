docker build . -f Dockerfile.GWS -t websocket_benchmark/go-gws && \
docker run --rm -p 8080:8080 websocket_benchmark/go-gws
