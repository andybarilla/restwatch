.DEFAULT_GOAL := dev

clean:
	rm -rf ./public/output.css
	rm -rf ./restwatch

build:
	npx @tailwindcss/cli -i input.css -o ./public/styles/app.css --minify
	GOARCH=wasm GOOS=js go build -o web/app.wasm cmd/wasm/main.go
	go build -o rest-watch cmd/frontend/main.go

run: build
	./rest-watch

tw:
	@npx @tailwindcss/cli -i input.css -o ./public/styles/app.css --watch
