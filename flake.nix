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
          (python3.pkgs.grip.overrideAttrs (old: {
            src = pkgs.fetchFromGitHub {
              owner = "joeyespo";
              repo = "grip";
              rev = "v4.6.1";
              sha256 = "sha256-CHL2dy0H/i0pLo653F7aUHFvZHTeZA6jC/rwn1KrEW4=";
            };

              patches = [ ];
          }))
        ];

      };

      packages.jutge = pkgs.buildGo118Module {
        pname = "jutge";
        version = "0.3.1";
        src = ./.;
        vendorSha256 = "sha256-xUwORIAWICnYOfApp8p5hBuaXwbzVVDOUtIPM9QATSI=";

        buildInputs = [ pkgs.installShellFiles ];

        postInstall = ''
          cat <<EOF >jutge.fish
          function __complete_jutge
              set -lx COMP_LINE (commandline -cp)
              test -z (commandline -ct)
              and set COMP_LINE "$COMP_LINE "
              $out/bin/jutge
          end
          complete -f -c jutge -a "(__complete_jutge)"
          EOF

          cat <<EOF >jutge.zsh
          autoload -U +X bashcompinit && bashcompinit
          complete -C $out/bin/jutge jutge
          EOF

          echo "complete -C $out/bin/jutge jutge" >jutge.bash

          installShellCompletion --cmd jutge jutge.{bash,fish,zsh}
        '';
      };

      defaultPackage = packages.jutge;
    }
  );
}
