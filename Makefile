build:
	@go build -o bin/portfolio-backend
run: build
	@./bin/portfolio-backend
test:
	@go test -v ./...
docker:
	@docker buildx build --tag kangym-go-backend .