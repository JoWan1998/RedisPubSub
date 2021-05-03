FROM golang
WORKDIR /usr/local/pub
COPY . .
RUN go mod download
CMD ["go","run","room.go"]
EXPOSE 8080