FROM alpine:latest AS BentoInstaller

ENV BENTO4_VERSION=1.6.0-637
ENV BENTO4_INSTALL_DIR=/opt/bento4
ENV PATH=/opt/bento4/bin:${PATH}

# Install dependencies.
RUN apk update \
  && apk add --no-cache \
  ca-certificates bash python3 wget libgcc make gcc g++

# Fetch source.
RUN cd /tmp/ \
  && wget -O Bento4-${BENTO4_VERSION}.tar.gz https://github.com/axiomatic-systems/Bento4/archive/v${BENTO4_VERSION}.tar.gz \
  && tar -xzvf Bento4-${BENTO4_VERSION}.tar.gz && rm Bento4-${BENTO4_VERSION}.tar.gz

# Create installation directories.
RUN mkdir -p \
  ${BENTO4_INSTALL_DIR}/bin \
  ${BENTO4_INSTALL_DIR}/scripts \
  ${BENTO4_INSTALL_DIR}/include

# Build.
RUN cd /tmp/Bento4-${BENTO4_VERSION}/Build/Targets/x86-unknown-linux \
  && make AP4_BUILD_CONFIG=Release

# Install.
RUN cd /tmp \
  && cp -r Bento4-${BENTO4_VERSION}/Build/Targets/x86-unknown-linux/Release/. ${BENTO4_INSTALL_DIR}/bin \
  && cp -r Bento4-${BENTO4_VERSION}/Source/Python/utils/. ${BENTO4_INSTALL_DIR}/utils \
  && cp -r Bento4-${BENTO4_VERSION}/Source/Python/wrappers/. ${BENTO4_INSTALL_DIR}/bin \
  && cp -r Bento4-${BENTO4_VERSION}/Source/C++/**/*.h . ${BENTO4_INSTALL_DIR}/include

# Cleanup.
RUN rm -rf /var/cache/apk/* /tmp/*

FROM golang:1.16-alpine

ENV PATH="$PATH:/bin/bash" \
    BENTO4_BIN="/opt/bento4/bin" \
    PATH="$PATH:/opt/bento4/bin"

RUN touch /root/.bashrc | echo "PS1='${debian_chroot:+($debian_chroot)}\u@\h:\w\$ '" >> /root/.bashrc

# FFMPEG + BASH + CURL
RUN apk add --update ffmpeg bash curl gcc g++ python3

# INSTALAÇÃO BENTO4
RUN mkdir -p /opt/bento4

# BUILD
COPY --from=BentoInstaller /opt/bento4 /opt/bento4

WORKDIR /go/src

#VOU MUDAR PARA O ENDPOINT CORRETO. USANDO TOP APENAS PARA SEGURAR O PROCESSO RODANDO
ENTRYPOINT [ "top" ]