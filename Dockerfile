FROM golang:1.25.4-alpine AS builder
WORKDIR /app
COPY go.mod ./
# RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o app .

# Final stage
FROM scratch

WORKDIR /app
COPY --from=builder /app/app /bin/app
COPY --from=builder /app/base /app/base
COPY --from=builder /app/bonus /app/bonus
EXPOSE 8080
ENTRYPOINT ["/bin/app"]
