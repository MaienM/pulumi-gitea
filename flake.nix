{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";

    flake-utils.url = "github:numtide/flake-utils";

    poetry2nix.url = "github:nix-community/poetry2nix";
    poetry2nix.inputs.nixpkgs.follows = "nixpkgs";
    poetry2nix.inputs.flake-utils.follows = "flake-utils";
  };

  outputs = { nixpkgs, flake-utils, ... }@inputs: flake-utils.lib.eachDefaultSystem (system:
    let
      pkgs = nixpkgs.legacyPackages.${system};
      poetry2nix = inputs.poetry2nix.lib.mkPoetry2Nix {
        inherit pkgs;
      };

      cruft = poetry2nix.mkPoetryApplication {
        projectDir = pkgs.fetchFromGitHub {
          owner = "Chilipp";
          repo = "cruft";
          rev = "d829bccec3bf690fe24387a1810795a05dd650c1";
          sha256 = "sha256-Up/r4e6jO5t5rTly05OpvHnuuyxxSJhNappq6kg51Zk=";
        };
        preferWheels = true;
        overrides = poetry2nix.overrides.withDefaults (_: super: {
          mkdocs-material = super.mkdocs-material.overridePythonAttrs (old: {
            postPatch = ''
              touch pyproject.toml
              ${old.postPatch or ""}
            '';
          });
        });
      };

      make = pkgs.writeShellApplication {
        name = "make";
        runtimeInputs = with pkgs; [ gnumake ];
        text = ''
          make "$@" SHELL=${pkgs.bash}/bin/bash
        '';
      };
    in
    {
      packages = {
        inherit cruft;
      };

      devShell = pkgs.mkShell {
        buildInputs = with pkgs; [

          # Base
          cruft
          make
          pulumi
          pulumictl

          # NodeJS
          nodejs
          nodejs.pkgs.yarn
          typescript

          # Python
          (python3.withPackages (pkgs: with pkgs; [
            packaging
            setuptools
          ]))

          # Go
          go
          golangci-lint
          gopls

          # Dotnet
          dotnet-sdk

          # Java
          (callPackage gradle-packages.gradle_8 {
            java = jdk11;
          })
        ];
      };
    }
  );
}
