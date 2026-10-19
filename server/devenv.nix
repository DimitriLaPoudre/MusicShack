{
  pkgs,
  lib,
  config,
  inputs,
  ...
}:

{
  # https://devenv.sh/basics/
  # env.GREET = "devenv";

  # https://devenv.sh/packages/
  packages = [
    pkgs.jq
  ];

  # https://devenv.sh/languages/
  languages.go = {
    enable = true;
    version = "1.27.1";
  };

  # https://devenv.sh/processes/
  processes.server = {
    exec = ''
      secretspec run -- go run ./cmd/server/main.go
    '';
    cwd = "${config.git.root}/server";
    restart = {
      on = "on_failure";
      max = 5;
    };
    ready = {
      http.get = {
        port = 8080;
        path = "/health";
        # host = "127.0.0.1"; # default
        # scheme = "http"; # default
      };
    };
    watch = {
      paths = [
        ./.
      ];
      extensions = [
        "go"
        "sum"
        "mod"
      ];
    };
  };

  # https://devenv.sh/services/
  # services.postgres.enable = true;

  # https://devenv.sh/scripts/
  # scripts.server.exec = ''
  # '';

  # https://devenv.sh/basics/
  enterShell = ''
    echo server devenv loaded !
  '';

  # https://devenv.sh/tasks/
  # tasks = {
  #   "myproj:setup".exec = "mytool build";
  #   "devenv:enterShell".after = [ "myproj:setup" ];
  # };

  # https://devenv.sh/tests/
  # enterTest = ''
  #   echo "Running tests"
  #   git --version | grep --color=auto "${pkgs.git.version}"
  # '';

  # https://devenv.sh/git-hooks/
  # git-hooks.hooks.shellcheck.enable = true;

  # See full reference at https://devenv.sh/reference/options/
}
