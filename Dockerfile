FROM golang:1.24.1
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o main .
EXPOSE 7540
CMD ["./main"]