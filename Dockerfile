FROM golang
WORKDIR /app
COPY ./ /app
RUN go build /app/cmd/main.go
CMD [ "/app/main", "--config=./downloader.yaml"]