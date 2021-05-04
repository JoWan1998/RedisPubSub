FROM golang
WORKDIR /usr/local/Pub
COPY . /usr/local/Pub
RUN go mod download
RUN go get -d -v ./...
RUN go install -v ./...
RUN go mod tidy
CMD ["go","run","room.go"]
EXPOSE 8080