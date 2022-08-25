package main

import (
	"fmt"
	"os"

	gfmt "github.com/PinkD/gofimpt/fmt"
)

func printHelpAndExit() {
	fmt.Printf(`Usage:
  %s [file/dir...]

    if no arg is provided, all files tracked by this git repo will formatted
    if files or directories are provided, provided files and all files under directories will be formatted
`, os.Args[0])
	os.Exit(1)
}

func parseFileAndModuleFromArg() ([]string, string) {
	var files []string
	dir := "."
	switch {
	case len(os.Args) == 1:
		// no args, detect git
		gitFiles, err := modifiedGitFiles(".")
		if err != nil {
			fmt.Println(err)
			printHelpAndExit()
		}
		files = append(files, gitFiles...)
	default:
		arg := os.Args[1]
		if arg == "-h" || arg == "--help" {
			printHelpAndExit()
		}
		filenames := os.Args[1:]
		for _, filename := range filenames {
			f, err := os.Stat(filename)
			if err != nil {
				fmt.Println(err)
				printHelpAndExit()
			}
			if f.IsDir() {
				dir = filename
				dirFiles, err := findGoFileUnderDir(filename)
				if err != nil {
					fmt.Println(err)
					printHelpAndExit()
				}
				files = append(files, dirFiles...)
			} else if isGoFile(filename) {
				files = append(files, filename)
			}
		}
	}
	module, err := getModuleName(dir)
	if err != nil {
		panic(err)
	}
	return files, module
}

func main() {
	files, module := parseFileAndModuleFromArg()
	errs := Run(files, func(_ int, file string) error {
		return gfmt.FormatFile(module, file)
	})
	for _, err := range errs {
		if err != nil {
			panic(err)
		}
	}
}
