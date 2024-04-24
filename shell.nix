{pkgs ? import <nixpkgs> {}}: let
  lib = pkgs.lib;
  # # minitime is a mini-output time wrapper.
  # minitime =
  #   pkgs.writeShellScriptBin "minitime"
  #   "command time --format $'%C -> %es\\n' \"$@\"";
in
  pkgs.mkShell {
    name = "apero-go";

    packages = with pkgs; [
      # gobject-introspection
      gtk3
      gtk-layer-shell
      # gtk4-layer-shell
      librsvg
      # pkg-config

      ### Build tools
      go
      gopls
      go-tools
      air

      just
      inotify-tools
      sassc

      # minitime
    ];

    # LD_PRELOAD = "${pkgs.gtk4-layer-shell}/lib/libgtk4-layer-shell.so";
    LD_PRELOAD = "${pkgs.gtk-layer-shell}/lib/libgtk-layer-shell.so";

    shellHook = ''
      export PATH="${pkgs.go}/bin:${pkgs.gopls}/bin:${pkgs.gotools}/bin:$PATH"
    '';
  }
