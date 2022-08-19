FROM golang:1.19.0-alpine3.16 as builder
RUN apk --no-cache add ca-certificates
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 go build -o /app

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app /app
EXPOSE 8080
CMD [ "/app" ]