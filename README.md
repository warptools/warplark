# warplark

Warplark is a tool for turning [Starlark](https://github.com/bazelbuild/starlark) code into inputs for [Warpforge](https://github.com/warptools/warpforge) inputs.

This is a **prototype** that is intended to be used as a generator within Warpforge.

## Installing

The easiest way to install `warplark` is using the Go CLI:

```
go install github.com/warptools/warplark@latest
```

Alternatively, clone this repository and install using Go:

```
git clone git@github.com:warptools/warplark
cd warplark
go install ./...
```

## Examples 

The `examples/` folder contains a variety of examples. These can be evaluted into Warpforge plots, then run with Warpforge. This will require that the `warplark` binary is in your `$PATH`.

```
warpforge -v plan generate examples/...
warpforge run examples/...
```
