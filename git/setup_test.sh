#!/bin/bash

set -eo pipefail -o nounset

function setup() {
     cd "$1"

     git init

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
