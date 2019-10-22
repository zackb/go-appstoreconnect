OUT=connect

default: build

build:
	go build -o $(OUT)

updatedeps:
	go list -m -u all

cleandeps:
	go mod tidy

clean:
	rm -f $(OUT)
