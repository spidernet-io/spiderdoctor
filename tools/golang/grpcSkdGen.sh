#!/bin/bash

# Copyright 2022 Authors of spidernet-io
# SPDX-License-Identifier: Apache-2.0

set -o errexit
set -o nounset
set -o pipefail

PROJECT_ROOT=$(dirname ${BASH_SOURCE[0]})/../..

# ====modify====
API_CODE_DIR=${1:-"${PROJECT_ROOT}/api/v1/agentGrpc"}

#======================

# CONST
PROTOC_GEN_GO_CODEGEN_PKG=${CODEGEN_PKG:-$(cd ${PROJECT_ROOT}; ls -d -1 ./vendor/google.golang.org/grpc/cmd/protoc-gen-go-grpc 2>/dev/null )}
PROTOC_GEN_CODEGEN_PKG=${CODEGEN_PKG:-$(cd ${PROJECT_ROOT}; ls -d -1 ./vendor/google.golang.org/protobuf/cmd/protoc-gen-go/ 2>/dev/null )}

PROTOC_GEN_Cmd() {
  go run ${PROJECT_ROOT}/${PROTOC_GEN_CODEGEN_PKG}/main.go "$@"
}

PROTOC_GO_GEN_Cmd() {
  go run ${PROJECT_ROOT}/${PROTOC_GEN_GO_CODEGEN_PKG}/main.go "$@"
}

