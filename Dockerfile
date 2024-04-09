FROM golang

RUN apt-get update && \
    apt-get install -y curl git tmux && \
    rm -rf /var/lib/apt/lists/*
    
ENV FOUNDRY_DIR=/
RUN curl -L https://foundry.paradigm.xyz | bash && \
    foundryup

COPY go.mod go.sum /app/
COPY go-ethereum/go.mod go-ethereum/go.sum /app/go-ethereum/
COPY dela/go.mod dela/go.sum /app/dela/
COPY kyber/go.mod kyber/go.sum /app/kyber/
WORKDIR /app

RUN go mod download
