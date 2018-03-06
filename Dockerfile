FROM iron/go

EXPOSE 8080

WORKDIR /app
ADD alchemy /app/

ENTRYPOINT ["./alchemy"]
