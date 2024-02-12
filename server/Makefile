.PHONY: run build-templ build-css build-go

run: build-templ build-css build-go
	@echo "Starting application..."
	@./tmp/app

build-templ:
	@echo "Generating Templ files..."
	@templ generate

build-css:
	@echo "Building Tailwind CSS..."
	@tailwindcss build --config ./tailwind.config.js -i ./static/tailwind.css -o ./static/styles.css

build-go:
	@echo "Building Go application..."
	@go build -o ./tmp/app ./app
