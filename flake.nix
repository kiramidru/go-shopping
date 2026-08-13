{
  description = "Go Flake";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/fbcf476f790d8a217c3eab4e12033dc4";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            air
          ];
        };
      }
    );
}
