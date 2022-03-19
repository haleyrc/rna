package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	DefaultOrg       = "haleyrc"
	DefaultGoVersion = "1.16"
)

func main() {
	if err := run(os.Args...); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func run(args ...string) error {
	return runCommand(args[1], args[2:]...)
}

func runCommand(cmd string, args ...string) error {
	cmd = strings.ToLower(cmd)
	switch cmd {
	case "new":
		return newProject(args[0], args[1:]...)
	default:
		return fmt.Errorf("invalid command: %s", cmd)
	}
}

func newProject(t string, args ...string) error {
	t = strings.ToLower(t)
	switch t {
	case "cli", "test":
		proj := CommandLineProject{
			TemplatePath: filepath.Join("templates", t),
		}

		if err := proj.ParseFlags(args...); err != nil {
			return fmt.Errorf("new project: %w", err)
		}

		if err := proj.Validate(); err != nil {
			return fmt.Errorf("new project: %w", err)
		}

		log.Printf("creating project: %s\n", proj.Name)

		if err := proj.Create(); err != nil {
			return fmt.Errorf("new project: %w", err)
		}

		if err := proj.PostCreate(); err != nil {
			return fmt.Errorf("new project: %w", err)
		}

		return nil
	default:
		return fmt.Errorf("invalid project type: %s", t)
	}
}

func makeRegularFile(path, newPath string) error {
	log.Printf("creating file: %s\n", newPath)
	if err := copyFile(path, newPath); err != nil {
		return fmt.Errorf("make regular file: %w", err)
	}
	return nil
}

func copyFile(src, dest string) error {
	fileBytes, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("copy file: %w", err)
	}

	if err := ioutil.WriteFile(dest, fileBytes, os.ModePerm); err != nil {
		return fmt.Errorf("copy file: %w", err)
	}

	return nil
}

func makeDirectory(path string) error {
	log.Printf("creating directory: %s/\n", path)
	if err := os.Mkdir(path, os.ModePerm); err != nil {
		return fmt.Errorf("make directory: %w", err)
	}
	return nil
}

func loadTemplate(path string) (string, error) {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("load template: %w", err)
	}

	return string(fileBytes), nil
}

func makeTemplatedFile(name, templatePath, outputPath string) error {
	log.Printf("creating file: %s\n", outputPath)
	tmpl, err := loadTemplate(templatePath)
	if err != nil {
		return fmt.Errorf("make templated file: %w", err)
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("make templated file: %w", err)
	}
	defer out.Close()

	t, err := template.New("path").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("make templated file: %w", err)
	}

	if err := t.Execute(out, struct {
		GoVersion string
		Name      string
		Org       string
	}{
		GoVersion: DefaultGoVersion,
		Name:      name,
		Org:       DefaultOrg,
	}); err != nil {
		return fmt.Errorf("make templated file: %w", err)
	}

	return nil
}
