# Goranger
A dropin replacement for the popular Python based Ranger utility. It should be
compatible with Ranger's Vim plugin.

## Installing
```
go install github.com/j18e/goranger@latest
```

## Using with vim
Use the ranger plugin at `github.com/francoiscabrol/ranger.vim` and in your
`.vimrc` place the following line:
```
let g:ranger_command_override = "$GOPATH/bin/goranger"
```

## Flags, behavior
The `--choosefiles` flag should take any highlighted files once they're entered
into and instead of opening them for editing, write their full paths to the
given file from the `--choosefiles` flag.

The `--selectfile` flag should start goranger while highlighting the file with
the given name, path not required.
