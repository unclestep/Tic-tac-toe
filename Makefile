swagger:
	swag init -g cmd/main.go

run: swagger
	CGO_ENABLED=0 go run cmd/main.go

.PHONY: swagger run
