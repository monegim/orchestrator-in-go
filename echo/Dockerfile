FROM golang:1.20-alpine3.18
WORKDIR /build
ADD *.go go.* ./
RUN go build .

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /build/echo .

CMD /app/echo