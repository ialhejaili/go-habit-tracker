FROM golang:1.22.3 AS builder
WORKDIR /app
COPY go.mod go.sum .env ./
RUN go mod download
COPY repository/ repository/ 
COPY model/ model/
COPY test/ test/ 
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o go-habit-tracker

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/go-habit-tracker ./
COPY --from=builder /app/.env ./
CMD ["./go-habit-tracker"]
