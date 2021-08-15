#!/bin/bash

if [[ ! -f ".env" ]]; then
    cp .env.example .env
fi

if [[ ! -d /go/src/tmp ]]; then
  mkdir /go/src/tmp
fi

go mod tidy

top