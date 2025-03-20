#!/bin/bash

# SPDX-FileCopyrightText: 2025 Red Hat, Inc. <sd-mt-sre@redhat.com>
#
# SPDX-License-Identifier: Apache-2.0

set -eo pipefail -o nounset

function setup() {
     cd "$1"

     git init

     git config --global user.name "test"
     git config --global user.email "test@example.com"

     touch file

     git add file

     git commit -m "initial commit"

     git tag "v1.0.0"

     echo "data" >> file

     git add file

     git commit -m "adding data"

     git tag "v2.0.0"

     git checkout -b "test"
}

function main() {
     setup "$1"
}

main "$@"
