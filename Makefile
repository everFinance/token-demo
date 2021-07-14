all:
	go build -o ./build/issuer ./cmd/issuer

test:
	sh ./test.sh