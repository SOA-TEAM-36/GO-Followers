FROM golang:alpine as build-stage
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o follower-service

FROM alpine
COPY --from=build-stage app/follower-service /usr/bin
EXPOSE 8084
ENTRYPOINT [ "follower-service" ]