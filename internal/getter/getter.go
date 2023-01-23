package getter

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	getter "github.com/hashicorp/go-getter"
	"golang.org/x/exp/slices"
)

func getWd() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return wd
}

func InitWorkdirPath(tmpWorkdir string) string {
	if _, err := os.Stat(tmpWorkdir); !os.IsNotExist(err) {
		err = os.RemoveAll(tmpWorkdir)
		if err != nil {
			log.Fatal(err)
		}
	}
	return tmpWorkdir
}

func InitWorkdir() string {
	wd := getWd()
	tmpWorkdir := filepath.Join(wd, ".temp_hcl2diff")
	return InitWorkdirPath(tmpWorkdir)
}

func GetFiles(srcToGet, workdir, target string, extensions []string) []string {
	targetWd := filepath.Join(workdir, target)
	get(srcToGet, targetWd)
	return findFiles(targetWd, extensions)
}

func findFiles(dirPath string, extensions []string) []string {
	var dir string
	symDir, err := os.Readlink(dirPath)
	if err == nil {
		dir = symDir
	} else {
		dir = dirPath
	}

	var files []string
	err = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && slices.Contains(extensions, filepath.Ext(path)) {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Failed to get files with extensions %v: %v\n", extensions, err)
	}

	return files
}

func get(src, dst string) {
	wd := getWd()
	client := &getter.Client{
		Ctx: context.Background(),
		//define the destination to where the directory will be stored. This will create the directory if it doesnt exist
		Dst: dst,
		Dir: true,
		Pwd: wd,
		//the repository with a subdirectory I would like to clone only
		Src:  src,
		Mode: getter.ClientModeDir,
		//define the type of detectors go getter should use, in this case only github is needed
		Detectors: []getter.Detector{
			&getter.GitHubDetector{},
			&getter.FileDetector{},
		},
		//provide the getter needed to download the files
		Getters: map[string]getter.Getter{
			"git":  &getter.GitGetter{},
			"file": &getter.FileGetter{},
		},
	}
	//download the files
	if err := client.Get(); err != nil {
		fmt.Fprintf(os.Stderr, "Error getting path %s: %v", client.Src, err)
		os.Exit(1)
	}
}
