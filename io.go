package main

import (
	"fmt"
	"github.com/krrrr38/gpshow/utils"
	"io/ioutil"
	"os"
	"sync"
)

// StaticDownload has external url and save point
type StaticDownload struct {
	url      string
	outDir   string
	filename string
}

// CopyResourceDir copy template files into project dir
func CopyResourceDir(srcDir, outDir string) {
	os.MkdirAll(outDir, 0755)
	targetResourcePath := "resources/" + srcDir
	templateFiles, err := AssetDir(targetResourcePath)
	utils.DieIf(err)
	for _, filename := range templateFiles {
		copyResourceDir(targetResourcePath, outDir, filename)
	}
}

func copyResourceDir(srcDir, outDir, filename string) {
	srcPath := srcDir + "/" + filename
	outPath := outDir + "/" + filename
	_, err := AssetInfo(srcPath)

	if err != nil {
		err := os.Mkdir(outPath, 0755)
		utils.DieIf(err)

		srcFiles, err := AssetDir(srcPath)
		utils.DieIf(err)
		for _, nextFilename := range srcFiles {
			copyResourceDir(srcPath, outPath, nextFilename)
		}
	} else {
		bytes, err := Asset(srcPath)
		utils.DieIf(err)
		ioutil.WriteFile(outPath, bytes, 0644)
	}
}

// DownloadStaticFiles try to download files and save them
func DownloadStaticFiles(targetDir string, files []StaticDownload) {
	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		go func(file StaticDownload) {
			bytes, err := FetchFile(file.url)
			if !utils.ErrorIf(err) {
				dirpath := targetDir + "/" + file.outDir
				filepath := dirpath + file.filename
				if _, err := os.Stat(dirpath); err != nil {
					os.MkdirAll(dirpath, 0755)
				}
				utils.ErrorIf(ioutil.WriteFile(filepath, bytes, 0644))
			}
			wg.Done()
		}(file)
	}
	wg.Wait()
}

// CopyLocalStaticFiles copy static dir into outdir
func CopyLocalStaticFiles(outDir, showPath string, ignoreDirs []string) {
	files, _ := ioutil.ReadDir(showPath)
	ignoreDirs = append(ignoreDirs, outDir)
	for _, file := range files {
		if file.Name() != "conf.js" && !arrayContain(file.Name(), ignoreDirs) {
			copyDir(showPath, outDir, file)
		}
	}
}

func arrayContain(target string, arr []string) bool {
	for _, str := range arr {
		if str == target {
			return true
		}
	}
	return false
}

func copyDir(srcDir, outDir string, file os.FileInfo) {
	outFilepath := fmt.Sprintf("%s/%s", outDir, file.Name())
	srcFilepath := fmt.Sprintf("%s/%s", srcDir, file.Name())
	if file.IsDir() {
		os.Mkdir(outFilepath, file.Mode())
		files, err := ioutil.ReadDir(srcFilepath)
		if !utils.ErrorIf(err) {
			for _, child := range files {
				copyDir(srcFilepath, outFilepath, child)
			}
		}
	} else {
		bytes, err := ioutil.ReadFile(srcFilepath)
		if !utils.ErrorIf(err) {
			ioutil.WriteFile(outFilepath, bytes, file.Mode())
		}
	}
}
