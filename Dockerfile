FROM docker.io/golang:latest

WORKDIR /app

COPY pogo.go .

# setup project and dependencies
RUN go mod init pogo
RUN go mod tidy

# prepare patched whatsmeow lib
RUN git clone --filter=blob:none https://github.com/tulir/whatsmeow.git
RUN go mod edit -replace go.mau.fi/whatsmeow=./whatsmeow

COPY add_prekey_timestamp.patch .
COPY disable_device_cache.patch .

RUN git -C whatsmeow/ apply ../add_prekey_timestamp.patch
RUN git -C whatsmeow/ apply ../disable_device_cache.patch

# build
RUN go build -v -o pogo


CMD ["/app/pogo"]
