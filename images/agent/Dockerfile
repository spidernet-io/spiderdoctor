# Copyright 2022 Authors of spidernet-io
# SPDX-License-Identifier: Apache-2.0

#====modify====
ARG BASE_IMAGE=ghcr.io/spidernet-io/spiderdoctor/agent-base:bfd810bfe3a2073060d813a6aa0b906207d4c6b6
ARG GOLANG_IMAGE=docker.io/library/golang:1.19@sha256:dc76ef03e54c34a00dcdca81e55c242d24b34d231637776c4bb5c1a8e8514253

# TARGETARCH is an automatic platform ARG enabled by Docker BuildKit.
# like amd64 arm64
ARG TARGETARCH

# report error by github ci
#ARG BUILDPLATFORM

#======= build bin ==========
FROM --platform=${BUILDPLATFORM} ${GOLANG_IMAGE} as builder

# docker build var
ARG TARGETOS
ARG TARGETARCH

# custom var
ARG RACE
ARG NOSTRIP
ARG NOOPT

COPY . /src
WORKDIR /src
RUN  make GOARCH=${TARGETARCH}   \
        RACE=${RACE} NOSTRIP=${NOSTRIP} NOOPT=${NOOPT} \
        DESTDIR_BIN=/tmp/install/${TARGETOS}/${TARGETARCH}/bin \
        build_agent_bin



#====== release image =======

FROM ${BASE_IMAGE}

LABEL maintainer="maintainer@spidernet-io"

# TARGETOS is an automatic platform ARG enabled by Docker BuildKit.
ARG TARGETOS
# TARGETARCH is an automatic platform ARG enabled by Docker BuildKit.
ARG TARGETARCH

ARG GIT_COMMIT_VERSION
ENV ENV_GIT_COMMIT_VERSION=${GIT_COMMIT_VERSION}
ARG GIT_COMMIT_TIME
ENV ENV_GIT_COMMIT_TIMESTAMP=${GIT_COMMIT_TIME}
ARG VERSION
ENV ENV_VERSION=${VERSION}

COPY --from=builder /tmp/install/${TARGETOS}/${TARGETARCH}/bin/*   /usr/bin/

CMD ["/usr/bin/agent"]