FROM golang:alpine as golang

WORKDIR /app

# Sao chép toàn bộ mã nguồn, bao gồm cả .env
COPY . .

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /server .

FROM scratch

# Sao chép binary từ stage build
COPY --from=golang /server .

# Sao chép file .env từ stage build
COPY --from=golang /app/.env /app/.env

# Đặt biến môi trường để ứng dụng biết vị trí file .env
ENV ENV_FILE_PATH=/app/.env

EXPOSE 8080

CMD ["/server"]
