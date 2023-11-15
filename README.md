# `nogo`
## minimalistic `go`-based `cli` app to interact with Notion

### usage

```shell
# list all commands
nogo
# or
nogo h
```

first you should configure the Notion API key (encrypted & stored locally in `$HOME/.config/nogo` by default) and the page ID you want to interact with.

```shell
# open configure prompt with
nogo c
```

#### `nogo` stack functionality
```shell
# show the stack
nogo s

# show the help message stack
nogo s -h
```

current commands:
```shell
NAME:
   nogo stack - show the stack

USAGE:
   nogo stack command [command options] [arguments...]

COMMANDS:
   add, a     add a new entry to the stack
   mod, m     modify a stack entry
   toggle, t  toggle stack entries
   rm, r      remove stack entries
   help, h    Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help (default: false)
```

## Dev

Publishing steps:

```shell
git commit -m '<COMMENT>'
git tag <VERSION>
git push origin <VERSION>
GOPROXY=proxy.golang.org go list -m github.com/haykh/nogo@<VERSION>
```
