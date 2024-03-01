BINARY_NAME = gac


default: build

build: 
	go build -o $(BINARY_NAME)

.PHONY: linux mac macintel windows

linux:
	GOOS=linux GOARCH=amd64 go build -o ${BINARY_NAME}-linux-amd64

mac:
	GOOS=darwin GOARCH=arm64 go build -o ${BINARY_NAME}-mac-arm64

macintel:
	GOOS=darwin GOARCH=amd64 go build -o ${BINARY_NAME}-mac-amd64

windows:
	GOOS=windows GOARCH=amd64 go build -o ${BINARY_NAME}-windows-amd64.exe

clean:
	rm -f ${BINARY_NAME}*