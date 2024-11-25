# FROM golang:1.18-alpine as builder

# RUN mkdir /app

# COPY . /app

# WORKDIR /app

# RUN CGO_ENABLED=0 go build -o accountLink ./cmd/main

# RUN chmod +x /app/accountLink

FROM alpine:latest

RUN mkdir /app

# COPY --from=builder /app/accountLink /app
COPY accountUpdate /app
COPY /db /db

CMD ["/app/accountUpdate"]
