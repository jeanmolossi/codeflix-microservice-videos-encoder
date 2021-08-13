FROM alpine:latest AS BentoInstaller

ENV BENTO4_BASE_URL="https://www.bok.net/Bento4/binaries" \
    BENTO4_VERSION="1-6-0-639" \
    BENTO4_CHECKSUM="05d6ed5c4254f1210f78c87b6003769a82d81500" \
    BENTO4_PATH="/opt/bento4" \
    BENTO4_TMP_PATH="/tmp/bento4"

ENV BENTO4_FILENAME="Bento4-SDK-$BENTO4_VERSION.x86_64-unknown-linux.zip" \
    BENTO4_BINARIES="$BENTO4_BASE_URL/Bento4-SDK-$BENTO4_VERSION.x86_64-unknown-linux.zip"

# Install Dependencies
RUN apk update && \
    apk add --no-cache curl unzip bash

RUN curl -O -s ${BENTO4_BINARIES} && \
    sha1sum -b ${BENTO4_FILENAME} | grep -o "^$BENTO4_CHECKSUM "  && \
    mkdir -p ${BENTO4_PATH} && \
    unzip ${BENTO4_FILENAME} -d ${BENTO4_TMP_PATH} && \
    rm -rf ${BENTO4_FILENAME} && \
    apk del unzip

# Install
RUN cd ${BENTO4_TMP_PATH}
RUN mv ${BENTO4_TMP_PATH}/Bento4-SDK-*.x86_64-unknown-linux/* ${BENTO4_PATH}

FROM golang:1.16-alpine

ENV PATH="$PATH:/bin/bash" \
    BENTO4_BIN="/opt/bento4/bin" \
    PATH="$PATH:/opt/bento4/bin"

# FFMPEG + BASH + CURL
RUN apk add --update ffmpeg bash curl

# INSTALAÇÃO BENTO4
WORKDIR /opt/bento4

# BUILD
COPY --from=BentoInstaller /opt/bento4 /opt/bento4

WORKDIR /go/src

#VOU MUDAR PARA O ENDPOINT CORRETO. USANDO TOP APENAS PARA SEGURAR O PROCESSO RODANDO
ENTRYPOINT [ "top" ]