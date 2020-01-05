build-linux: 
	GOOS=linux GOARCH=amd64 go build -o build/rsec-net

run: build-linux
	docker-compose up
