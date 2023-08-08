package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"go.starlark.net/lib/json"
	"go.starlark.net/starlark"
)

type Pragmas struct {
	version *int
}
type loader struct {
	basePath string
}

// Load will start a starlark thread for a given module.
// Warplark interprets "modules" to be filepaths.
// Filepaths are scanned from the loader basepath+module name up to the root directory.
func (l *loader) Load(_ *starlark.Thread, module string) (starlark.StringDict, error) {
	// find the module file by checking parent directories
	// until the first match is found
	file := filepath.Join(l.basePath, module)
	for {
		_, err := os.Stat(file)
		if err == nil {
			// found it!
			break
		}
		if file == "/"+module {
			// hit root directory, fail
			panic(fmt.Sprintf("failed to locate module %q", module))
		}
		if os.IsNotExist(err) {
			// file not found, check the parent dir
			currentDir := filepath.Dir(file)
			parentDir := filepath.Dir(currentDir)
			file = filepath.Join(parentDir, module)
			continue
		}
		// error doing stat, fail
		panic(fmt.Sprintf("error locating module %q: %s", module, err))
	}
	absFile, err := filepath.Abs(file)
	if err != nil {
		panic(err)
	}
	newLoader := loader{basePath: filepath.Dir(absFile)}
	thread := &starlark.Thread{Name: "module " + module, Load: newLoader.Load}
	globals, err := starlark.ExecFile(thread, file, nil, nil)

	return globals, err
}

func parsePragmas(path string) (Pragmas, error) {
	txt, err := ioutil.ReadFile(path)
	if err != nil {
		return Pragmas{}, err
	}

	// this regex matches a string in the format of
	// #+ warplark pragma-key pragma-value
	pragmaRe, err := regexp.Compile("#\\+warplark\\s+([^\\s]+)\\s+(.*)")
	if err != nil {
		return Pragmas{}, err
	}

	pragmas := Pragmas{}

	lines := strings.Split(string(txt), "\n")
	for lineNum, line := range lines {
		line = strings.TrimSpace(line)
		matches := pragmaRe.FindStringSubmatch(line)
		if matches == nil {
			// no pragma on this line, stop parsing pragmas
			break
		}
		pragmaKey := strings.ToLower(matches[1])
		pragmaValue := matches[2]

		switch pragmaKey {
		case "version":
			version, err := strconv.Atoi(pragmaValue)
			if err != nil {
				return Pragmas{}, fmt.Errorf("failed to parse version pragma. expected int, got %q", pragmaValue)
			}
			pragmas.version = &version
			continue
		}
		return Pragmas{}, fmt.Errorf("unknown pragma %q at %s:%d", pragmaKey, path, lineNum)
	}

	return pragmas, nil
}

// Execute Starlark program in a file.
func execFile(path string) error {

	pragmas, err := parsePragmas(path)
	if err != nil {
		panic(err)
	}

	if pragmas.version == nil {
		panic("no version pragma defined!")
	}
	if *pragmas.version != 0 {
		panic("unsupported warplark version")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	newLoader := loader{basePath: filepath.Dir(absPath)}
	thread := &starlark.Thread{Name: "my thread", Load: newLoader.Load}
	globals, err := starlark.ExecFile(thread, path, nil, starlark.StringDict{"json": json.Module})
	if err != nil {
		switch serr := err.(type) {
		case *starlark.EvalError:
			log.Println(serr.Backtrace())
			return err
		}
		log.Println("unrecognized starlark error type")
		return err
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
	return nil
}

func main() {
	log.SetFlags(0)
	if len(os.Args) == 2 {
		err := execFile(os.Args[1])
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		return
	}
	log.Printf("usage: %s [input file]\n", os.Args[0])
}
