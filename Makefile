BINARY_NAME=ego

build:
	go build -o ${BINARY_NAME} ./main.go

clean:
	rm ${BINARY_NAME}
