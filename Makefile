build:
	@go build -o bin/portfolio-backend
run: build
	@./bin/portfolio-backend
test:
	@go test -v ./...