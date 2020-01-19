.PHONY: build

HOST_BUILD=build/rsec-net

build:
	go build -o $(HOST_BUILD)

build-linux:
	GOOS=linux GOARCH=amd64 go build -o build/linux/rsec-net

run: build-linux
	docker-compose up --build

run-single: build-linux
	docker-compose -f docker-compose.debug.yml up --build

run-host: build
	./$(HOST_BUILD) start-node -v \
	--nodeName=debug \
	--announceMulticast=false \
	--announceAddr=localhost:1100 \
	--announceListenPort=1140 \
	--dataListenPort=1148

debug: build-linux run-host

clean:
	docker-compose down
