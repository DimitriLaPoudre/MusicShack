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
    pkgs.git

    # backend
    pkgs.air
    pkgs.jq

    # frontend
  ];

  # https://devenv.sh/languages/
  languages = {
    go = {
      enable = true;
      lsp.enable = true;
    };

    javascript = {
      enable = true;
      directory = "./client_web";
      lsp.enable = true;
      bun = {
        enable = true;
        install.enable = true;
      };
    };
  };

  # https://devenv.sh/processes/
  # processes.dev.exec = "${lib.getExe pkgs.watchexec} -n -- ls -la";

  # https://devenv.sh/services/
  # services.postgres.enable = true;

  # https://devenv.sh/scripts/
  scripts.front-dev.exec = ''
    	cd "$DEVENV_ROOT"/client_web
    	secretspec run --provider dotenv:../.env -- bun run dev
  '';

  scripts.back-dev.exec = ''
    	cd "$DEVENV_ROOT"/server
    	secretspec run --provider dotenv:../.env -- air | jq
  '';

  # https://devenv.sh/basics/
  # enterShell = ''
  #   hello         # Run scripts directly
  #   git --version # Use packages
  # '';

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
