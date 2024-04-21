FROM golang:1.21.1 AS build-stage
    WORKDIR /api 

    COPY go.mod go.sum ./
    RUN go mod download 

    COPY .*.go ./
    RUN CGO_ENABLED=0 GOOS=linux go build -o /api

FROM build-stage AS run-test-stage 
    RUN go test -v ./...

FROM scratch AS run-release-stage 
