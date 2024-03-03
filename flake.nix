{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    #nixpkgs.url = "github:NixOS/nixpkgs?ref=39cca54ab0f547b5a59959f9cf4541ea63a53220";

    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, }:

  flake-utils.lib.eachDefaultSystem (system:
    let
      pkgs = import nixpkgs {
        inherit system;
        overlays = [(final: prev: {
          go = prev.go_1_22;
          buildGoModule = prev.buildGo122Module;
          buildGoPackage = prev.buildGo122Package;
        })];
      };
    in
    rec {
      devShell = pkgs.mkShellNoCC {
        name = "go";

        buildInputs = with pkgs; [
          go
          gopls
          gotools
        ];
      };

      packages.hello = pkgs.buildGo118Module {
        pname = "hello";
        version = "0.0.1";
        src = ./.;
        vendorSha256 = "";
      };

      defaultPackage = packages.hello;
    }
  );
}
