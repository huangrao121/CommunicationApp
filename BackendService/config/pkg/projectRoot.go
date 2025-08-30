package pkg

import (
	"log"
	"os"
	"path/filepath"
)

func GetProjectRoot() string {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exPath := filepath.Dir(ex)
	log.Printf("Project root: %s", exPath)

	// 向上查找 go.mod 文件
	for {
		if _, err := os.Stat(filepath.Join(exPath, "go.mod")); err == nil {
			return exPath
		}
		parent := filepath.Dir(exPath)
		if parent == exPath {
			log.Fatal("go.mod not found, please run this program in the root of the project")
			break
		}
		exPath = parent
	}
	wd, _ := os.Getwd()
	return wd
}
