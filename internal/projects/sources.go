package projects

import (
	"os"
	"path/filepath"
)

func LoadAllProjects(directories []string) ([]string, error) {
	var projects []string

	for _, dir := range directories {
		subdirectories, err := loadSubdirectories(dir)

		if err != nil {
			return nil, err
		}

		projects = append(projects, subdirectories...)
	}

	return projects, nil
}

func loadSubdirectories(dir string) ([]string, error) {
	var subdirs []string

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			subdirs = append(subdirs, filepath.Join(dir, file.Name()))
		}
	}

	return subdirs, nil
}
