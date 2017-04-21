# Sub-Alias
Apply alias to sub-commands

## Description  
I don't like longer commands. So I often use `alias`.  
Just like following:  
``` sh
$ alias d=docker
$ d ps
// docker ps
```

Infrequently, we want to use `alias` to sub-command.  
However, `alias` command can apply for command only.  

`salias` means sub-alias. `salias` makes it possible to apply alias to sub-commands.  
Therefore, for example, `docker inspect` can alias to `docker i`.  

## Equipments
- Go v1.8 or newer
- bash or zsh

## Installation
``` sh
$ go get github.com/lycoris0731/salias
```

## Usage
Add following command to `.bashrc` or `.zshrc`.  
``` sh
source <(salias __init__)
```

## License
Please see [LICENSE](./LICENSE).
