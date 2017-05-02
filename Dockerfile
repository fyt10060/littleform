FROM golang
RUN mkdir /app
ADD littleform /app/littleform
ADD conf /app/conf
WORKDIR /app
EXPOSE 9090
EXPOSE 6379
ENTRYPOINT /app/littleform

CMD ["bee", "run"]