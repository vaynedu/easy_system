Bin="marketing"

default: build lint

build:
	go build -o $(Bin)

lint:
	golint ./...

hz:
	pwd
	hz new -force -idl idl/hello.thrift
run:
	go build -o $(Bin) && ./$(Bin)