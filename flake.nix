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
              rustfmt
            ];
            
            preBuild = ''
              # Create bin directory
              mkdir -p bin
              
              # Build TypeScript backend
              cd backends/typescript
              npm install
              npm run build
              cd ../..
              
              # Build Rust backend
              cd backends/rust
              cargo build --release
              cp target/release/rust-rm-comments ../../bin/
              cd ../..
            '';
            
            postInstall = ''
              # Install language backends
              mkdir -p $out/libexec/rm-comments
              
              # Copy TypeScript backend
              mkdir -p $out/libexec/rm-comments/typescript
              cp -r backends/typescript/dist/* $out/libexec/rm-comments/typescript/
              
              # Copy Python backend  
              cp backends/python/python_rm_comments.py $out/libexec/rm-comments/
              chmod +x $out/libexec/rm-comments/python_rm_comments.py
              
              # Copy Rust backend
              if [ -f bin/rust-rm-comments ]; then
                cp bin/rust-rm-comments $out/libexec/rm-comments/
                chmod +x $out/libexec/rm-comments/rust-rm-comments
              fi
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
            rustfmt
            nixpkgs-fmt
          ];
        };
      });
}

