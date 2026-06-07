swagger:
	swag init -g cmd/main.go

run: swagger
	go run cmd/main.go

.PHONE: swagger run
