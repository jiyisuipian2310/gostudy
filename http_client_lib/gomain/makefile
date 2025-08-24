target=libgocommonlib.so

#nm -D libgocommonlib.so | grep ' T ' | grep -v "cgo"

all:
	#go build -v -x -buildmode=c-shared -o ${target}
	go build -buildmode=c-shared -o ${target}
	chmod +x ${target}

clean:
	rm -rf $(basename $(target))*
