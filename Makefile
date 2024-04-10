build: clean
	go build main.go -o apero

start:
	go run -v ./main.go start

fmt:
	fmt *.go

clean:
	rm -rf apero

godoc:
	godoc -http=:6060