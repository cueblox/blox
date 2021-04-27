# Development

## Required Tools

* Go > 1.16
* gofumpt `GO111MODULE=on go get mvdan.cc/gofumpt`
* [golangci-lint](https://golangci-lint.run/usage/install/#local-installation)


## Building

```shell
make setup # download go dependencies
make build # or just make
```

## Testing

```shell
make test
```
