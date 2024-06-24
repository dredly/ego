BINARY_NAME=ego
TEST_BINARY_NAME=${BINARY_NAME}_test

test:
	go test -timeout 30s -v ./internal/...

e2e:
	go build -o ${TEST_BINARY_NAME} ./main.go
	go test -timeout 30s -v github.com/dredly/ego/e2e

build:
	go build -o ${BINARY_NAME} ./main.go

clean:
	rm -f ${BINARY_NAME}
	rm -f ${TEST_BINARY_NAME}

.PHONY: test e2e build clean
