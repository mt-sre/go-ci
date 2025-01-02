#!/bin/bash

set -exvo pipefail -o nounset

# utilize local go 1.23 version if available
GO_1_23="/opt/go/1.23.1/bin"

if [ -d  "${GO_1_23}" ]; then
     PATH="${GO_1_23}:${PATH}"
fi

PROJECT_ROOT=$(git rev-parse --show-toplevel)
DEPENDENCY_BIN="${PROJECT_ROOT}/.cache/dependencies/bin"
PRECOMMIT="${DEPENDENCY_BIN}/pre-commit"

function install_precommit() {
     if [ -f "${PRECOMMIT}" ]; then
          echo "pre-commit already installed"

          return
     fi

     version=$1
     url="https://github.com/pre-commit/pre-commit/releases/download/v${version}/pre-commit-${version}.pyz"

     setup_dependencies && curl -sL "${url}" -o "${PRECOMMIT}" && chmod +x "${PRECOMMIT}"
}

function setup_dependencies() {
     mkdir -p "${DEPENDENCY_BIN}"
}

function run-hooks() {
     ${PRECOMMIT} run \
     --from-ref origin/main \
     --to-ref HEAD \
     --show-diff-on-failure
}

install_precommit "2.17.0" && run-hooks && go test -v -cover -race -count=1 ./...
