FROM golang:1.21-alpine
WORKDIR /src
COPY . .
RUN go mod download
RUN go build -o /bin/cms ./cmd/cms
ENTRYPOINT [ "/bin/cms" ]