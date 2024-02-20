# Specifies a parent image
FROM golang:1.20-alpine

RUN apk add cmd:bash cmd:tmux cmd:xxd
 
# Creates an app directory to hold your app’s source code
WORKDIR /app
 
# Copies everything from your root directory into /app
COPY ./ /app/
 
# Installs Go dependencies
RUN go install ./dkg/pedersen_bn256/dkgcli
