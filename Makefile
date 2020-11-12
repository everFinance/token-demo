all:
	go build -o ./build/issuer ./cmd/issuer
	go build -o ./build/detector ./cmd/detector

test:
	sh ./test.sh