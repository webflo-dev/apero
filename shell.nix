{pkgs ? import <nixpkgs> {}}: let
  lib = pkgs.lib;

  # minitime is a mini-output time wrapper.
  minitime =
    pkgs.writeShellScriptBin "minitime"
    "command time --format $'%C -> %es\\n' \"$@\"";
in
  pkgs.mkShell {
    name = "apero-go";

    packages = with pkgs; [
      gobject-introspection
      gtk4
      gtk4-layer-shell
      atk

      ### Build tools
      pkg-config
      go
      gopls
      go-tools
      # air

      minitime
    ];

    LD_PRELOAD = "${pkgs.gtk4-layer-shell}/lib/libgtk4-layer-shell.so";

    shellHook = ''
      export PATH="${pkgs.go}/bin:${pkgs.gopls}/bin:${pkgs.gotools}/bin:$PATH"
    '';
  }
