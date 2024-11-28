package projects

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

type Project struct {
	Name       string
	Path       string
	Running    bool
	Active     bool
	LastActive int
}

func LoadAllProjects(directories []string) (*Projects, error) {
	projects := Projects{
		Map: make(map[string]*Project),
	}

	for _, dir := range directories {
		absolutePath := transformPath(dir)
		projectList, err := loadProjectsInDirectory(absolutePath)

		if err != nil {
			return nil, err
		}

		for _, p := range projectList {
			projects.Map[p.Name] = &p
		}
	}

	return &projects, nil
}

func loadProjectsInDirectory(dir string) ([]Project, error) {
	var projects []Project

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() && file.Name()[0] != '.' {
			p := Project{
				Name:       file.Name(),
				Path:       filepath.Join(dir, file.Name()),
				LastActive: 0,
			}
			projects = append(projects, p)
		}
	}

	return projects, nil
}

func transformPath(dir string) string {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalf("cannot get home dir: %s\n", err)
	}

	return strings.ReplaceAll(dir, "~", home)
}
