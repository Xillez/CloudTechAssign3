FROM alpine:latest

MAINTAINER Edward Muller <edward@heroku.com>

WORKDIR "/opt"

ADD .docker_build/cmd/web/web /opt/bin/web
ADD .docker_build/cmd/clock/clock /opt/bin/web

CMD ["/opt/bin/web"]
CMD ["/opt/bin/clock"]
