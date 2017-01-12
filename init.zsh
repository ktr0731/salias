#!/usr/bin/env zsh

config="$HOME"/.salias

__salias::register() {
  if ! grep "$1" $config > /dev/null; then
    echo "$1" >> $config
  fi
}

salias() {
  if [ "$#" -ne 2 ]; then
    echo "invalid arguments"
  fi

  __salias::register "$1" "$2"
}
