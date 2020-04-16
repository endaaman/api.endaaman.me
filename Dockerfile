# FROM golang:alpine AS builder
FROM golang:buster

# system
RUN apt-get update
RUN apt-get install -y supervisor nginx-extras sudo

RUN \
  chown -R www-data:www-data /var/lib/nginx && \
  echo "\ndaemon off;" >> /etc/nginx/nginx.conf && \
  rm /etc/nginx/sites-enabled/default

# app
ENV GOPATH /var/www
ENV PATH $PATH:$GOPATH/bin
ENV APP_DIR "$GOPATH/src/github.com/endaaman/api.endaaman.me"
ENV BEEGO_RUNMODE prod

RUN mkdir -p $APP_DIR
WORKDIR $APP_DIR
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY nginx/api.conf /etc/nginx/sites-enabled/
COPY supervisor.conf /etc/supervisor/conf.d/
COPY . .
RUN go generate main.go
RUN go build -ldflags '-w -s' -o app

RUN mkdir -p /var/www/shared
RUN chmod 2775 /var/www
RUN chown -R www-data:www-data /var/www

VOLUME /var/www/shared
EXPOSE 80 443
CMD ["/usr/bin/supervisord"]
