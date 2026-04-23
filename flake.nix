{
  description = "phunter — terminal UI for hunting and killing processes on TCP ports";

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
        version = "0.3.0";
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "phunter";
          inherit version;
          src = ./.;

          # Run `nix build` once with the placeholder below — it will fail and
          # print the correct hash. Replace this value with that hash.
          vendorHash = "sha256-B+jvh+oSkIsna1zber6CdU5CX9vSricq+zlc7KzWoXM=";

          env.CGO_ENABLED = "0";

          ldflags = [
            "-s"
            "-w"
            "-X main.version=${version}"
          ];

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
