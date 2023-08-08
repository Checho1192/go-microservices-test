FROM alpine:latest

RUN mkdir /app

COPY statsMsApp /app

CMD ["/app/statsMsApp"]