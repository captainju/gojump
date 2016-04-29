# gojump
simple stupid autojump-like, tailored for fish shell

## Script to add in config.fish
```shell
begin
  set --local GOJUMP_PATH /path/to/gojump.fish
  if test -e $GOJUMP_PATH
    source $GOJUMP_PATH
  end
end```
