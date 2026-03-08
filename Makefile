OUT=connect

default: build

build:
	go build -o $(OUT) ./cmd/connect.go

test:
	go test -v github.com/zackb/go-appstoreconnect/appstoreconnect

updatedeps:
	go list -m -u all

cleandeps:
	go mod tidy

clean:
	rm -f $(OUT)
