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
``` bash
$ salias go i=install
$ salias docker c=container
$ exec $SHELL # reload the shell

$ go i github.com/golang/go
# `go install github.com/golang/go` 

$ docker c ls
# `docker container ls`
```

## Equipments
- Go v1.8 or newer
- bash, Zsh or fish

## Installation
``` sh
$ export GOPATH=~/go # if $GOPATH is not set
$ go get github.com/lycoris0731/salias
$ sudo ln $GOPATH/bin/salias /usr/bin/
```

## Usage
### Set sub-alias definition file
Please create one of these files:  
- $SALIAS_PATH

- $XDG_CONFIG_HOME/salias/salias.toml

- $HOME/salias.toml

Then use `salias <command> <subalias>=<subcommand>` to set sub-alias. Or edit `salias.toml` manually.

### Make sub-alias initializes automatically

Add following command.  

`.bashrc` or `.zshrc`.  

``` sh
source <(salias --init) # or `salias -i` for short
```

`config.fish`

``` sh
source (salias --init | psub)
```

## How It Works
When initialization, `salias` registers the command as salias's alias.  
``` sh
# [go]
# b = "build"
$ source <(salias --init)
$ type go
go is an alias for salias --run go
```

`salias` find sub-alias that is sub-command of passed command as arguments.  
If hit sub-alias, execute it.  
Or not found, execute as it is.  

## License
Please see [LICENSE](./LICENSE).
