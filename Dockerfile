# FROM golang:alpine AS builder
FROM golang:buster

# system
RUN apt-get update
RUN apt-get install -y supervisor nginx-extras sudo

RUN \
  chown -R www-data:www-data /var/lib/nginx && \
  echo "\ndaemon off;" >> /etc/nginx/nginx.conf && \
  rm /etc/nginx/sites-enabled/default

COPY nginx/nginx.conf /etc/nginx/sites-enabled
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
RUN go build -o app

# shared
RUN mkdir -p /var/www/shared
RUN chmod 2775 /var/www/shared
RUN chown -R www-data:www-data /var/www/shared

VOLUME /var/www/shared
EXPOSE 80 443
CMD ["/usr/bin/supervisord"]



# CMD /go/src/api/app -runmode=prod

# FROM alpine
# COPY --from=builder /go/src/api/app /app
# CMD /app -runmode=prod

# RUN go get github.com/beego/bee
# CMD bee run -downdoc=true -gendoc=true -runmode=prod
