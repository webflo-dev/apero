build-ctl:
	cd apero-ctl && go build -o aperoctl ./main.go

build-demo:
	cd demo && go build -o demo ./main.go

run-demo:
	cd demo && go run -v ./main.go

build: build-ctl build-demo

clean-build:
	rm -f apero-ctl/aperoctl
	rm -f demo/demo