# ghf



## Description

## Usage

Write pull-request format file like [sample_pullreq_template.md](https://github.com/deme0607/ghf/blob/master/sample_pullreq_template.md).
You can also use [`text/template`](https://golang.org/pkg/text/template/) style format file.

And specify your template file on `.ghf-tempalte` file at your project directory.

.ghf-tempalte:
```
sample_pullreq_template.md
```

After configuration, you can make pull request

```
$ ghf
```

command.
If you don't want to open editor, you can use `-no-editor` option.

For additional usage, please type `ghf -h` for help.

## Install

To install, use `go get`:

```bash
$ go get github.com/deme0607/ghf
```

or download binary from [release page](https://github.com/deme0607/ghf/releases) and place it in `$PATH` directory

## Requirement

[hub](https://github.com/github/hub) command is required.

## Contribution

1. Fork ([https://github.com/deme0607/ghf/fork](https://github.com/deme0607/ghf/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[deme0607](https://github.com/deme0607)
