run:
	go run cmd/api/main.go

build-a:
	docker build -t websocket_server:latest -f deploy/api/linux/Dockerfile .

run-a:
	docker run --name websocket_server -p 5488:5488 --network="host" -v ${PWD}/docker/log:/app/log websocket_server:latest


# docker exec -it <container-id> sh