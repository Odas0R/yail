#!/usr/bin/env bash

set -e

version=$(curl https://go.dev/VERSION?m=text)
echo "VERSION: $version"

os="linux"
arch="amd64"
tarFile="${version}.${os}-${arch}.tar.gz"
url="https://golang.org/dl/${tarFile}"

# Download the latest golang
wget --continue --show-progress "${url}"
printf "Downloaded Golang!\n"

# Remove the old golang
sudo rm -rf /usr/local/go

# Install the new Go
sudo tar -C /usr/local -xzf "$tarFile"
printf "Create the skeleton for your local users go directory\n"
mkdir -p ~/go
mkdir -p ~/go/bin
mkdir -p ~/go/pkg
mkdir -p ~/go/src

# GoPath
echo "export GOPATH=~/go" >>~/.bashrc
echo 'export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin' >>~/.bashrc

# Remove Download
sudo rm "$tarFile"

# Print Go Version
/usr/local/go/bin/go version

# Reload bashrc
. ~/.bashrc
