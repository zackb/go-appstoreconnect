OUT=connect

default: build

build:
	go build -o $(OUT)

test:
	go test -v github.com/zackb/go-appstoreconnect/appstoreconnect

updatedeps:
	go list -m -u all

cleandeps:
	go mod tidy

clean:
	rm -f $(OUT)
