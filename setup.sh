#!/bin/bash

go mod init pogo
go mod tidy

git clone --filter=blob:none https://github.com/tulir/whatsmeow.git

go mod edit -replace go.mau.fi/whatsmeow=./whatsmeow

git -C whatsmeow/ apply ../add_prekey_timestamp.patch
git -C whatsmeow/ apply ../disable_device_cache.patch
