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

## Writing Warplark Code

Warplark uses `starlark-go` to implement Starlark. The [language spec](https://github.com/google/starlark-go/blob/master/doc/spec.md) provides the full details of the language.

### Pragmas

Warplark files begin with special lines known as *pragmas*. These provide information about the file which Warplark uses to determine how to parse it. A pragma is written in the format:

```
#+warplark [pragma-name] [pragma-value]
```

Pragmas must be provided at the start of the file. Once a non-pragma line is hit, Warplark will stop processing pragams.

The `version` pragma is mandatory for all Warplark files. Version 0 is currently used to denote the alpha status of Warplark.

#### Supported Pragmas

| Pragma Name | Value Type | Description                      | Example                |
|-------------|------------|----------------------------------|------------------------|
| version     | int        | Version of Warplark used by file | `#+warplark version 0` |
