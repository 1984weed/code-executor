build:
	GOOS=freebsd GOARCH=386 go build -o bin/code-executor-freebsd-386 main.go
	GOOS=linux GOARCH=386 go build -o bin/code-executor main.go

run:
	go run main.go


docker-build:
	@docker build -t baracode/python -f ./docker-images/python/Dockerfile .  