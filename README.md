# dockerhand
A command line tool Docker resources manager

Please make sure you use dry-run before running any command.

Cleaning orphaned volumes is **EXPERIMENTAL** and should be used only if you are at the same machine Docker is running. Use at your own risk.

## Usage
    dockerhand [command]

    Available Commands:
      clean       Clean resources
      show        Show resources
      help        Help about any command

    Flags:
      -d, --dryrun=false: enable dryrun mode
      -e, --endpoint="unix:///var/run/docker.sock": docker endpoint uri
      -h, --help=false: help for dockerhand


    Use "dockerhand help [command]" for more information about a command.

Binary release
--------------
Download from [releases page](https://github.com/edgard/dockerhand/releases)
