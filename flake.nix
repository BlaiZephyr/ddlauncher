{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages = {
          default = pkgs.callPackage ./. { };
        };

        devShells = {
          default = pkgs.mkShell {
            inherit (self.packages.${system}.default) nativeBuildInputs;

            packages =
              with pkgs;
              [
                go
                golangci-lint
              ]
              ++ self.packages.${system}.default.buildInputs;
          };
        };
      }
    );
}
