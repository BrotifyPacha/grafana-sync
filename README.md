![build](https://github.com/brotifypacha/grafana-sync/actions/workflows/ci.yaml/badge.svg)
![coverage](https://img.shields.io/badge/coverage-69.6%25-yellowgreen)

# Grafana-sync

`grafana-sync` is a command-line tool for downloading Grafana dashboards as
folder structure.

## Features

- Download dashboards and folders from Grafana as JSON files
- Store dashboards and folders on local disk
- Print Grafana structure tree to stdout

## Installation

`grafana-sync` can be installed with Go:

```
go install github.com/brotifypacha/grafana-sync/cmd/grafana-sync@latest
```

## Usage

```
Usage: grafana-sync command [command-flags...]

Commands:

    print - prints grafana structure tree to stdout

        -host (required)
            grafana HTTP api host
        -recursive (default: true)
            whether or not it should print tree structure recursivly
        -folder-uid (default: root folder)
            grafana folder UID to use for this command


    download - downloads grafana structure to specified dir

        -host (required)
            grafana HTTP api host
        -path (required)
            local path to directory where dashboards requests will be stored
        -folder-uid (default: root folder)
            grafana folder UID to use for this command

```

### Examples

Print Grafana structure tree:

```
grafana-sync print -host http://grafana.example.com
```

Download dashboards and folders to the `./grafana-data` directory:

```
grafana-sync download -host http://grafana.example.com -path ./grafana-data
```

## Roadmap

- [x] Write readme.md
- [ ] Automate coverage shield/badge in README
- [ ] Restructure `./cmd/grafana-sync/main.go` for better extensibility
- [ ] Add support for more querie types


## Contributing

Contributions are welcome! If you find a bug or have an idea for a new feature,
please [open an issue](https://github.com/brotifypacha/grafana-sync/issues/new).
If you want to contribute code, please open a pull request. Before submitting a
pull request, make sure the tests pass and the code is properly formatted:

```
make test
make fmt
```

## License

`grafana-sync` is licensed under the [MIT License](https://opensource.org/license/mit/).

