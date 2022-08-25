package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/go-git/go-git/v5"

	"gofimpt/errors"
)

func Run(ts []string, f func(int, string) error) []error {
	n := runtime.NumCPU() * 2
	if len(ts) < n {
		n = len(ts)
	}
	var wg sync.WaitGroup
	n = 1
	errs := make([]error, len(ts))
	for i := 0; i < n; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j, t := range ts {
				if j%n == i {
					errs[j] = f(j, t)
				}
			}
		}()
	}
	wg.Wait()
	return errs
}

func isGoFile(file string) bool {
	return strings.HasSuffix(file, ".go")
}

func findGoFileUnderDir(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if isGoFile(path) {
				files = append(files, path)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func getModuleName(dir string) (string, error) {
	name := "go.mod"
	dir, err := filepath.Abs(dir)
	if err != nil {
		return "", errors.Trace(err)
	}
	_, err = os.Stat(filepath.Join(dir, name))
	for os.IsNotExist(err) {
		last := dir
		dir = filepath.Dir(dir)
		if last == dir {
			// file not found
			return "", nil
		}
		file := filepath.Join(dir, name)
		_, err = os.Stat(file)
	}
	data, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return "", errors.Trace(err, "open go.mod")
	}
	lines := strings.SplitN(string(data), "\n", 2)
	if len(lines) == 0 {
		return "", errors.New("module name not found in go.mod")
	}
	module := strings.TrimSpace(strings.TrimPrefix(lines[0], "module "))
	return module, nil
}

func modifiedGitFiles(path string) ([]string, error) {
	opt := &git.PlainOpenOptions{
		DetectDotGit: true,
	}
	repo, err := git.PlainOpenWithOptions(path, opt)
	if err != nil {
		return nil, errors.Trace(err, "open git repo on %s", path)
	}
	tree, err := repo.Worktree()
	if err != nil {
		return nil, errors.Trace(err, "open git work tree on %s", path)
	}
	statusMap, err := tree.Status()
	if err != nil {
		return nil, errors.Trace(err, "get git status on %s", path)
	}
	var files []string
	for name, status := range statusMap {
		if status.Staging != git.Untracked {
			if isGoFile(name) {
				files = append(files, name)
			}
		}
	}
	return files, nil
}
