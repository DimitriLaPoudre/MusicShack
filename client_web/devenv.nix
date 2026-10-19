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
  # packages = [ pkgs.git ];

  # https://devenv.sh/languages/
  languages = {
    javascript = {
      enable = true;
      directory = "./client_web";
      bun = {
        enable = true;
        install.enable = true;
      };
    };
  };

  # https://devenv.sh/processes/
  processes.client_web.exec = "client_web";

  # https://devenv.sh/services/
  # services.postgres.enable = true;

  # https://devenv.sh/scripts/
  scripts.client_web.exec = ''
    cd "$DEVENV_ROOT"/client_web
    secretspec run -- bun run dev
  '';

  # https://devenv.sh/basics/
  enterShell = ''
    echo client_web devenv loaded !
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
