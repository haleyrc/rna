package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type CommandLineProject struct {
	Name         string
	TemplatePath string
}

func (clp *CommandLineProject) Create() error {
	if err := filepath.WalkDir(clp.TemplatePath, clp.create); err != nil {
		return fmt.Errorf("new project: %w", err)
	}
	return nil
}

func (clp *CommandLineProject) ParseFlags(args ...string) error {
	fs := flag.NewFlagSet("cli", flag.ContinueOnError)
	fs.StringVar(&clp.Name, "name", "", "The name of the project")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}

	return nil
}

func (clp *CommandLineProject) PostCreate() error {
	if err := clp.build(); err != nil {
		return fmt.Errorf("post create: %w", err)
	}
	if err := clp.test(); err != nil {
		return fmt.Errorf("post create: %w", err)
	}
	return nil
}

func (clp *CommandLineProject) Validate() error {
	if clp.Name == "" {
		return fmt.Errorf("project name is required")
	}
	return nil
}

func (clp *CommandLineProject) build() error {
	target := filepath.Join("build", clp.Name)

	cmd := exec.Command("go", "build", "-o", target, ".")
	cmd.Dir = clp.Name

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build: %w", err)
	}

	return nil
}

func (clp *CommandLineProject) create(path string, de fs.DirEntry, err error) error {
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

func (clp *CommandLineProject) test() error {
	cmd := exec.Command("go", "test", "-v", "-count=1", ".")
	cmd.Dir = clp.Name
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("test: %w", err)
	}

	return nil
}
