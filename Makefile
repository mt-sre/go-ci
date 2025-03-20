# SPDX-FileCopyrightText: 2025 Red Hat, Inc. <sd-mt-sre@redhat.com>
#
# SPDX-License-Identifier: Apache-2.0

check: lint test
.PHONY: check

lint:
	pre-commit run \
		--show-diff-on-failure \
		--from-ref "origin/main" \
		--to-ref "HEAD"
.PHONY: lint

golangci-lint:
	golangci-lint run -v --new-from-rev HEAD --fix
.PHONY: golangci-lint

tidy:
	go mod tidy
.PHONY: tidy

verify:
	go mod verify
.PHONY: verify

check-flake:
	nix flake check --all-systems
.PHONY: check-flake

test: test-units

test-units:
	go test -v -race -count=1 -v ./...
.PHONY: test-units

COPYRIGHT ?= Red Hat, Inc. <sd-mt-sre@redhat.com>
LICENSE ?= Apache-2.0

reuse-apply:
	reuse annotate --copyright="${COPYRIGHT}" --license="${LICENSE}" --skip-unrecognised --skip-existing ./**/*
.PHONY: reuse-apply
