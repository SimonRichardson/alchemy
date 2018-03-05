FROM iron/go

EXPOSE 8080

WORKDIR /app
ADD discourse /app/

ENTRYPOINT ["./discourse"]
