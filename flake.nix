# SPDX-FileCopyrightText: 2025 Red Hat, Inc. <sd-mt-sre@redhat.com>
#
# SPDX-License-Identifier: Apache-2.0

{
  inputs.flake-utils.url = "github:numtide/flake-utils";
  inputs.nixpkgs.url = github:NixOS/nixpkgs/nixos-unstable;
  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
    let
        pkgs = import nixpkgs { inherit system; };
        devDeps = with pkgs; [
          git
          go_1_23
          pre-commit
          golangci-lint
          reuse
        ];
    in
    {
      devShells.default = pkgs.mkShell {
        buildInputs = devDeps;
      };
  });
}
