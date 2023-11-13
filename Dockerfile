FROM --platform=linux/amd64 golang:1.17 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./out/dist

FROM --platform=linux/amd64 alpine:latest

WORKDIR /app

COPY --from=builder /app/out/dist .

# Create a non-root user
RUN adduser -D -g '' appuser

# Change to the non-root user
USER appuser

ENV PORT=8080

EXPOSE 8080

CMD ["./dist"]
