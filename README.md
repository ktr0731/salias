# Sub-Alias
[![wercker status](https://app.wercker.com/status/33e127b80f4ea96bc02dc4bfac4dbdac/s/master "wercker status")](https://app.wercker.com/project/byKey/33e127b80f4ea96bc02dc4bfac4dbdac)  
Apply alias to sub-commands

![DEMO](https://user-images.githubusercontent.com/12953836/41504734-b98bbc96-7233-11e8-9435-8a6ebdecb587.gif)

## Description  
I don't like longer commands. So I often use `alias`.  
Just like following:  
``` sh
$ alias d=docker
$ d ps
# `docker ps`
```

Infrequently, we want to use `alias` to sub-command.  
However, `alias` command can apply for command only.  

`salias` means sub-alias. `salias` makes it possible to apply alias to sub-commands.  

## Example
~/salias.toml
``` toml 
[go]
i = "install"
b = "build"
r = "run"

[docker]
i = "image"
c = "container"

[docker-compose]
l = "logs -f"
```

``` bash
$ go i github.com/golang/go
# `go install github.com/golang/go` 

$ docker i ls
# `docker image ls`

$ alias d=docker
$ d c ls
# `docker container ls`
```

## Equipments
- Go v1.8 or newer
- bash, Zsh or fish

## Installation
``` sh
$ go get github.com/lycoris0731/salias
```

## Usage
### Set sub-alias definition file
Please set the file to one of following.  
- $SALIAS_PATH
- $XDG_CONFIG_HOME/salias/salias.toml
- $HOME/salias.toml

Add following command.  

`.bashrc` or `.zshrc`.  
``` sh
source <(salias __init__)
```

`config.fish`
``` sh
source (salias __init__ | psub)
```

## How It Works
When initialization, `salias` registers the command as salias's alias.  
``` sh
# [go]
# b = "build"
$ source <(salias __init__)
$ type go
go is an alias for salias go
```

`salias` find sub-alias that is sub-command of passed command as arguments.  
If hit sub-alias, execute it.  
Or not found, execute as it is.  

## License
Please see [LICENSE](./LICENSE).
