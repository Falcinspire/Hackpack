FROM alpine
WORKDIR /app
COPY service /app
EXPOSE 80
ENTRYPOINT ["./service"]
