{
  description = "Multi-language comment stripping utility";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages = {
          default = self.packages.${system}.rm-comments;
          
          rm-comments = pkgs.buildGoModule {
            pname = "rm-comments";
            version = "0.1.0";
            src = ./.;
            vendorHash = null;
            
            nativeBuildInputs = with pkgs; [
              nodejs
              python3
              rustc
              cargo
            ];
            
            buildInputs = with pkgs; [
              nodejs
              python3
              rustc
            ];
            
            preBuild = ''
              # Build TypeScript backend
              cd backends/typescript
              ${pkgs.nodejs}/bin/npm install
              ${pkgs.nodejs}/bin/npm run build
              cd ../..
              
              # Build Rust backend
              cd backends/rust
              ${pkgs.cargo}/bin/cargo build --release
              cp target/release/rust-rm-comments ../../bin/
              cd ../..
            '';
            
            postInstall = ''
              # Install language backends
              mkdir -p $out/libexec/rm-comments
              
              # Copy TypeScript backend
              cp -r backends/typescript/dist $out/libexec/rm-comments/typescript
              cp backends/typescript/package.json $out/libexec/rm-comments/typescript/
              
              # Copy Python backend  
              cp backends/python/python_rm_comments.py $out/libexec/rm-comments/
              
              # Copy Rust backend
              cp bin/rust-rm-comments $out/libexec/rm-comments/ || true
            '';
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            nodejs
            python3
            rustc
            cargo
            nixpkgs-fmt
          ];
        };
      });
}
