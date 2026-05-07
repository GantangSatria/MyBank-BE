from golang:1.25-alpine AS builder

WORKDIR /app

# install dependencies
RUN apk add --no-cache git

# copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# copy source code
COPY . .

# build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

# Stage 2 — run
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache tzdata

# set timezone
ENV TZ=Asia/Jakarta

# copy binary dari builder
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]