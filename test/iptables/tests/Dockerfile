# Compile tests in a golang container.
FROM golang:1.12 AS golang
ADD *.go /
RUN go build -o /tests /*.go

# Create the container to run tests in.
FROM ubuntu:16.04
RUN apt-get update && apt-get install -y iptables
COPY --from=golang /tests /tests
ENTRYPOINT ["/tests"]
