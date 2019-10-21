
OUT=connect

default: build

build:
	go build -o $(OUT)

updatedeps:
	go list -m -u all

clean:
	rm $(OUT)
