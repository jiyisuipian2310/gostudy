target=libgohttplib.so

all:
	go build -buildmode=c-shared -o ${target} main.go
	chmod +x ${target}

clean:
	rm -rf ${target}
