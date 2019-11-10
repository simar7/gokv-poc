clean:
	rm -f gokv-poc
	docker rm -f localdynamodb localredis

build:
	go build -o gokv-poc main.go

demo:
	@echo "building binary..."
	make build
	@echo "starting docker db containers..."
	make dynamodb
	@echo "writing to dbs..."
	./gokv-poc
	@echo "cleaning up..."
	make clean

bench:
	@echo "performing benchmarks..."
	./gokv-poc -op=bench

dynamodb:
	docker run -p 8000:8000 --name localdynamodb -d amazon/dynamodb-local

all:
	make clean
	make build