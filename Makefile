swagger:
	swag init -g cmd/main.go

run: swagger
	go run cmd/main.go CGO_ENABLED=0

.PHONY: swagger run
