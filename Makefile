.DEFAULT_GOAL := dev

htmx:
	wget -O public/scripts/htmx.min.js https://cdn.jsdelivr.net/npm/htmx.org@2.0.5/dist/htmx.min.js
	wget -O public/scripts/htmx-ext-sse.min.js https://cdn.jsdelivr.net/npm/htmx-ext-sse@2.2.2

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
