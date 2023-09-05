# `nogo`
## Minimalistic `go`-based `cli` app to interact with Notion

### Usage

```shell
# list all commands
nogo
# or
nogo h
```

First you should configure the Notion API key (encrypted & stored locally in `$HOME/.config/nogo` by default) and the page ID-s of Notion pages you want to interact with.

```shell
# open configure prompt with
nogo c
```

#### `nogo` stack functionality
```shell
# show the stack help (todo list)
nogo s

# list the stack
nogo s -l

# add a new to-do item & list the stack
nogo s -la "new todo item"
nogo s -a "another todo item"

# remove a to-do item that contains "new"
nogo s -r "new"

# mark a to-do item as done (vague search)
nogo s -d "anot"
# mark as not done & list
nogo s -lu "anot"
```

## To-do

- [x] minimal structure of `cli`
- [x] notion api
  - [x] key storage/encryption 
  - [x] submitting/reading
