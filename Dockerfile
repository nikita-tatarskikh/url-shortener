FROM alpine

RUN apk --no-cache add ca-certificates

COPY url_shortener /url_shortener/

WORKDIR /url_shortener

CMD ["/bin/sh", "-c", "./url_shortener"]
