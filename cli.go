package main

import (
	"flag"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

type CommandLineProject struct {
	Name         string
	TemplatePath string
}

func (clp *CommandLineProject) Build() error {
	if err := filepath.WalkDir(clp.TemplatePath, clp.build); err != nil {
		return fmt.Errorf("new project: %w", err)
	}
	return nil
}

func (clp *CommandLineProject) build(path string, de fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	newPath := filepath.Join(
		clp.Name,
		strings.TrimPrefix(path, clp.TemplatePath),
	)

	if de.IsDir() {
		return makeDirectory(newPath)
	}

	if strings.HasSuffix(path, ".tmpl") {
		newPath = strings.TrimSuffix(newPath, ".tmpl")
		return makeTemplatedFile(clp.Name, path, newPath)
	}

	return makeRegularFile(path, newPath)
}

func (clp *CommandLineProject) ParseFlags(args ...string) error {
	fs := flag.NewFlagSet("cli", flag.ContinueOnError)
	fs.StringVar(&clp.Name, "name", "", "The name of the project")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}

	return nil
}

func (clp *CommandLineProject) Validate() error {
	if clp.Name == "" {
		return fmt.Errorf("project name is required")
	}
	return nil
}
