#################################
# STEP 1 build executable binary
#################################
FROM golang:1.16.3-alpine3.13 AS builder
WORKDIR $GOPATH/src/github.com/llamadeus/keyval-server

# Fetch dependencies
COPY go.mod go.sum ./
# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download -x

# Build the binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -o /go/bin/keyval-server

#############################
# STEP 2 build a small image
#############################
FROM alpine:3.13.5

# Copy our static executable
COPY --from=builder /go/bin/keyval-server /usr/bin/keyval-server

# Run the keyval-server binary
ENTRYPOINT ["/usr/bin/keyval-server", "-s", "/var/keyval-server/data.json"]
