FROM golang:alpine as golang

WORKDIR /app

COPY . .

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /server .

FROM scratch

COPY --from=golang /server .

# COPY --from=golang /app/.env /app/.env

# ENV ENV_FILE_PATH=/app/.env

EXPOSE 8080

CMD ["/server"]
