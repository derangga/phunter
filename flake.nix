{
  description = "phunter — terminal UI for hunting and killing processes on TCP ports";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in {
        packages.default = pkgs.buildGoModule {
          pname = "phunter";
          version = "0.1.0";
          src = ./.;

          # Run `nix build` once with the placeholder below — it will fail and
          # print the correct hash. Replace this value with that hash.
          vendorHash = "sha256-FwfpQvOVHmeS4KQuMhOUX/Hc4GDunG+fsT2ZMwLIJUo=";

          env.CGO_ENABLED = "0";

          meta = with pkgs.lib; {
            description = "Terminal UI for hunting and killing processes listening on TCP ports";
            homepage = "https://github.com/derangga/phunter";
            license = licenses.mit;
            mainProgram = "phunter";
          };
        };

        apps.default = {
          type = "app";
          program = "${self.packages.${system}.default}/bin/phunter";
        };
      }
    );
}
