git clone --filter=blob:none https://github.com/tulir/whatsmeow.git

go mod edit -replace go.mau.fi/whatsmeow=./whatsmeow

git -C whatsmeow/ apply ../add_prekey_timestamp.patch
