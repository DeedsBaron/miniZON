FROM ubuntu:22.04

ADD ./bin/app /app
COPY ./config.yaml /
ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

CMD ["/app"]
