FROM golang:1.19.1-alpine3.16 as build

RUN apk --no-cache add ca-certificates

COPY . /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o wisdom .

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/wisdom /

ENTRYPOINT [ "/wisdom" ]