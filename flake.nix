{
  description = "Go-based tool to do awesome stuff with notion";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
  };

  outputs =
    { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
    in
    {
      packages.${system}.default = pkgs.buildGoModule rec {
        pname = "nogo";
        version = "1.6.0";
        author = "haykh";

        src = pkgs.fetchFromGitHub {
          owner = author;
          repo = pname;
          rev = "v${version}";
          hash = "sha256-ItkpkZIhA9a6QgxFVSNN9YGQoML6NT4uoiu9aLsZI9o=";
        };

        vendorHash = "sha256-9aIcu1BImY7+IdNEVb3acgdM3kBamKrWVOyUnaSZXZk=";

        meta = with pkgs.lib; {
          description = "go-based tool to do awesome stuff with notion";
          homepage = "https://github.com/${author}/nogo";
          license = licenses.bsd3;
          maintainers = [ author ];
        };
      };

      devShells.${system}.default = pkgs.mkShell {
        packages = [
          pkgs.go
          pkgs.gopls
          pkgs.gotools
        ];

        shellHook = ''
          echo "Welcome to nogo development shell!"
        '';
      };
    };
}
