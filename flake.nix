{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    systems.url = "github:nix-systems/default";
    devenv.url = "github:cachix/devenv/1e4701fb1f51f8e6fe3b0318fc2b80aed0761914";
  };

  nixConfig = {
    extra-trusted-public-keys = "devenv.cachix.org-1:w1cLUi8dv3hnoSPGAuibQv+f9TZLr6cv/Hm9XgU50cw=";
    extra-substituters = "https://devenv.cachix.org";
  };

  outputs = { self, nixpkgs, devenv, systems, ... } @ inputs:
    let
      forEachSystem = nixpkgs.lib.genAttrs (import systems);
    in
    {
      devShells = forEachSystem
        (system:
          let
            pkgs = nixpkgs.legacyPackages.${system};
            googleapis = pkgs.fetchFromGitHub {
              owner = "googleapis";
              repo = "googleapis";
              rev = "e0677a395947c2f3f3411d7202a6868a7b069a41";
              hash = "sha256-9XRc2fOTV+yS6NcPpjcoW8e1F+jsaOyhu6zCUUvM0O4=";
            };
          in
          {
            default = devenv.lib.mkShell {
              inherit inputs pkgs;
              modules = [
                {
                  languages.go = {
                    enable = true;
                  };

                  packages = with pkgs; [ goose flyctl ];

                  enterShell = ''
                    echo "nBits chaintips shell activated!"

                  '';
                }
              ];
            };
          });
    };
}
