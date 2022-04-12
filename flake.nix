{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

    flake-utils = {
      url = "github:numtide/flake-utils";
      inputs.nixpkgs.follows = "nixpkgs";
    };

  };

  outputs = { self, nixpkgs, flake-utils, }:

  flake-utils.lib.eachDefaultSystem (system:
    let
      pkgs = import nixpkgs { inherit system; };
    in
    rec {
      devShell = pkgs.mkShellNoCC {

        name = "go_1.18";

        buildInputs = with pkgs; [
          go_1_18
          gopls
          gotools
        ];

      };

      packages.jutge = pkgs.buildGo118Module {
        pname = "jutge";
        version = "0.3.1";
        src = ./.;
        vendorSha256 = "sha256-swYBwxYa+mj9hKp0D/Wt06TMIywFut75rYQQ3E0NhEc=";

        buildInputs = [ pkgs.installShellFiles ];

        postInstall = ''
          installShellCompletion --cmd jutge \
            --bash <($out/bin/jutge --completion-script-bash) \
            --zsh <($out/bin/jutge --completion-script-zsh)
        '';
      };

      defaultPackage = packages.jutge;
    }
  );
}
