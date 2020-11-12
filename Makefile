all:
	go build -o ./build/detector ./cmd/detector

test:
	sh ./test.sh