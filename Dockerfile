FROM alpine
WORKDIR /app
COPY build/service /app
EXPOSE 80
ENTRYPOINT ["./service"]