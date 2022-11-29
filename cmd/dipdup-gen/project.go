package main

import (
	"os"
	"path/filepath"
)

type projectDirs struct {
	root     string
	storage  string
	postgres string
	cmd      string
	build    string
}

func createProjectDirs(rootDir, app string) (*projectDirs, error) {
	if rootDir == "" {
		output, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		rootDir = output
	}

	dirs := new(projectDirs)
	dirs.root = rootDir

	dirs.storage = filepath.Join(rootDir, "internal", "storage")
	if err := os.MkdirAll(dirs.storage, os.ModePerm); err != nil {
		if !os.IsExist(err) {
			return nil, err
		}
	}
	dirs.postgres = filepath.Join(dirs.storage, "postgres")
	if err := os.MkdirAll(dirs.postgres, os.ModePerm); err != nil {
		if !os.IsExist(err) {
			return nil, err
		}
	}
	dirs.cmd = filepath.Join(rootDir, "cmd", app)
	if err := os.MkdirAll(dirs.cmd, os.ModePerm); err != nil {
		if !os.IsExist(err) {
			return nil, err
		}
	}
	dirs.build = filepath.Join(rootDir, "build")
	if err := os.MkdirAll(dirs.build, os.ModePerm); err != nil {
		if !os.IsExist(err) {
			return nil, err
		}
	}

	return dirs, nil
}
