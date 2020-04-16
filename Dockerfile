# FROM golang:alpine AS builder
FROM golang:buster

# system
RUN apt-get update
RUN apt-get install -y supervisor nginx-extras sudo

RUN \
  chown -R www-data:www-data /var/lib/nginx && \
  echo "\ndaemon off;" >> /etc/nginx/nginx.conf && \
  rm /etc/nginx/sites-enabled/default

COPY nginx/nginx.conf /etc/nginx/sites-enabled/
COPY supervisor.conf /etc/supervisor/conf.d/

# app
RUN mkdir -p /var/www/app
ENV GOPATH /var/www/
ENV PATH $PATH:$GOPATH/bin

WORKDIR /var/www/app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags '-w -s' -o app

# shared
RUN mkdir -p /var/www/shared
RUN chmod 2775 /var/www/shared
RUN chown -R www-data:www-data /var/www/shared

ENV BEEGO_RUNMODE prod
VOLUME /var/www/shared
EXPOSE 80 443
CMD ["/usr/bin/supervisord"]
