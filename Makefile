build-linux: 
	GOOS=linux GOARCH=amd64 go build -o build/rsec-net

run: build-linux
	docker-compose up --build

debug: build-linux
	docker-compose -f docker-compose.debug.yml up --build

clean:
	docker-compose down
