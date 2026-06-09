FROM golang:latest AS build-stage
WORKDIR /app/
COPY go.mod go.sum ./
RUN go mod download && \
    go install github.com/swaggo/swag/cmd/swag@latest
COPY . .
RUN swag init -g cmd/main.go && \
    CGO_ENABLED=0 go build -o server ./cmd/main.go

FROM scratch
COPY --from=build-stage /app/server /server
COPY --from=build-stage /app/docs /docs/
COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
USER 1001
EXPOSE ${TICTACTOE_PORT}
CMD ["/server"]
