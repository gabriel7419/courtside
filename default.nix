{pkgs ? import <nixpkgs> {}, ...}:
pkgs.buildGoModule {
  pname = "golazo";
  version = "0.21.0";
  vendorHash = "sha256-M2gfqU5rOfuiVSZnH/Dr8OVmDhyU2jYkgW7RuIUTd+E=";

  subPackages = ["."];

  src = builtins.path {
    path = ./.;
    name = "source";
  };
}
