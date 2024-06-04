#!/bin/bash

sudo pacman -Syu --noconfirm go
go mod init mymodule
go get golang.org/x/crypto@latest

mkdir ~/.pswdcnt/module/
mv module/store.go ~/.pswdcnt/module/

go build ~/.pswdcnt/module/store.go
rm -r ~/.pswdcnt/module/store.go

sudo mv bash/wifi.sh /usr/local/bin/wifi
sudo chmod +x /usr/local/bin/wifi

