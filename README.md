# Useful utils for Go command-line programs

### `cli` package

* `cli.Args` struct: a better `os.Args`
* `cli.Register()`, `cli.Lookup()`: register and lookup CLI subcommands

### `api` package

* `api.Client` struct: thin wrapper around `net/http`

### `pathname` package

* `pathname.Pathname` struct: Ruby `Pathname`-like handling of filesystem paths

### `utils` package

* `utils.Env` struct: a map-based approach to `os.Environ()`
* `utils.Set` struct
