package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"go.starlark.net/lib/json"
	"go.starlark.net/starlark"
)

var basePath string

func load(_ *starlark.Thread, module string) (starlark.StringDict, error) {
	// find the module file by checking parent directories
	// until the first match is found
	file := filepath.Join(basePath, module)
	for {
		_, err := os.Stat(file)
		if err == nil {
			// found it!
			break
		} else if file == "/"+module {
			// hit root directory, fail
			panic(fmt.Sprintf("failed to locate module %q", module))
		} else if os.IsNotExist(err) {
			// file not found, check the parent dir
			currentDir := filepath.Dir(file)
			parentDir := filepath.Dir(currentDir)
			file = filepath.Join(parentDir, module)
		} else {
			// error doing stat, fail
			panic(fmt.Sprintf("error locating module %q: %s", module, err))
		}
	}

	thread := &starlark.Thread{Name: "module " + module, Load: load}
	globals, err := starlark.ExecFile(thread, file, nil, nil)

	return globals, err
}

// Execute Starlark program in a file.
func execFile(path string) {
	basePath = filepath.Dir(path)
	thread := &starlark.Thread{Name: "my thread", Load: load}
	globals, err := starlark.ExecFile(thread, path, nil, starlark.StringDict{"json": json.Module})
	if err != nil {
		panic(err)
	}

	// retrieve the starlark json functions
	json_encode := json.Module.Members["encode"]
	json_indent := json.Module.Members["indent"]

	// retrieve the plot to output
	plot := globals["result"]
	plotv1 := starlark.NewDict(1)
	plotv1.SetKey(starlark.String("plot.v1"), plot)

	// json encode then indent the plot using starlark library
	v, err := starlark.Call(thread, json_encode, starlark.Tuple{plotv1}, nil)
	if err != nil {
		panic(err)
	}
	v, err = starlark.Call(thread, json_indent, starlark.Tuple{v}, nil)
	if err != nil {
		panic(err)
	}

	// unescape the resulting string
	plotJson, err := strconv.Unquote(v.String())
	if err != nil {
		panic(err)
	}

	// print the result
	fmt.Println(plotJson)
}

func main() {
	if len(os.Args) == 2 {
		execFile(os.Args[1])
	} else {
		fmt.Printf("usage: %s [input file]\n", os.Args[0])
	}
}
