package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var dag map[string]map[string]int
func main() {
	// config
	pkgRootName := *flag.String("app", "", "package name")
	path := *flag.String("path", "", "app path")
	dotFilePath := *flag.String("dot-filepath", "./dag.dot", "dot out path")
	flag.Parse()
	path, _ = filepath.Abs(path)
	// check
	if !Exist(path) {
		fmt.Print("app path can not use")
		return
	}
	if pkgRootName == "" {
		fmt.Println("pkgRootName can not be nil")
		return
	}
	// start
	fmt.Printf("packageName: %s, path: %s, dot filepath: %s \n", pkgRootName, path, dotFilePath)
	fmt.Println("start...")
	dag = make(map[string]map[string]int)
	parse(pkgRootName, path)
	for k, vm := range dag {
		for vk, _ := range vm {
			fmt.Println(k, vk)
		}
	}
	writeToDotFile(dotFilePath)
	fmt.Println("end...")
}


func parse(pkgRootName, path string) {
	err := filepath.Walk(path, func(singlePath string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n",singlePath , err)
			return err
		}
		// skip vendor
		if info.IsDir() && info.Name() == "vendor"{
			//fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
			return filepath.SkipDir
		}
		// select go file
		if !info.IsDir() && strings.Contains(info.Name(), ".go") {
			parseSingleFile(pkgRootName, path, singlePath)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", path, err)
	}
}

func parseSingleFile(pkgRootName, path, singlePath string) {
	fmt.Println("start...", singlePath)
	f, err := os.Open(singlePath)
	if err != nil {
		panic(fmt.Sprintf("parse single file %s error:%s", singlePath, err.Error()))
	}
	r := bufio.NewReader(f)
	isImport := false
	packageName := ""
	for {
		line,_, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(fmt.Sprintf("read single file %s error:%s", singlePath, err.Error()))
		}
		if strings.HasPrefix(string(line), "package") {
			packageNames := strings.Split(string(line), " ")
			packageName += filepath.Join(packageName, packageNames[len(packageNames) -1])
		}
		if string(line) == "import (" {
			isImport = true
			continue
		}
		if isImport == true && string(line) == ")"{
			isImport = false
			break
		}
		if isImport == true {
			importName := importPkgName(string(line), pkgRootName)
			if importName == "" {
				continue
			}
			if _, ok := dag[packageName]; ok {
				dag[packageName][importName] = 1
			}else {
				tmp := make(map[string]int)
				tmp[importName] = 1
				dag[packageName] = tmp
			}
			fmt.Println("dag:", packageName, importName)
		}
	}
}

func writeToDotFile(dotFilePath string) {
	fd, _ := os.OpenFile(dotFilePath, os.O_RDWR|os.O_CREATE, 0644)
	defer fd.Close()
	fd.Write([]byte("digraph G {\n"))
	for k, vm := range dag {
		for v, _ := range vm {
			fd.Write([]byte(fmt.Sprintf("\t\"%s\" -> \"%s\"\n", k, v)))
		}
	}
	fd.Write([]byte("}\n"))
}

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}


func importPkgName(line, pkgRootName string) string{
	if index := strings.Index(line, pkgRootName); index != -1 {
		name := strings.Trim(line[index:], "\"")
		names := strings.Split(name, "/")
		return names[len(names) - 1]
	}
	return ""

}

