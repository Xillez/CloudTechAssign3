FROM alpine:latest

MAINTAINER Edward Muller <edward@heroku.com>

WORKDIR "/opt"

ADD .docker_build/Assignment4 /opt/bin/Assignment4

CMD ["/opt/bin/Assignment4"]
