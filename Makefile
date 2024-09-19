build: build-css build-go

build-css:
	templ generate
	npm run build:css

build-go:
	go build -o build/main ./src

run: build
	./build/main

clean:
	rm -f build/main
	rm -rf build/css/*

serve: ./build/main

dev:
	air