FROM golang
WORKDIR /pub
COPY . /pub
RUN go mod download
RUN go get -d -v ./...
RUN go install -v ./...
RUN go mod tidy
EXPOSE 80
CMD ["go","run","room.go"]
