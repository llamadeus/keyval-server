# keyval-server

A super basic key-value-storage with ttl and file persistence.

## Usage

Run the server using:
```shell
go run main.go -s storage/data.json
```

Or you can build the server and run the binary directly:
```shell
go build
./keyval-server -s storage/data.json
```

## Docker

```shell
# Building the image
docker build -t keyval-server .

# Run the container
docker run --rm --volume "$(pwd)"/storage:/var/keyval-server -p 3000:3000 keyval-server
```
