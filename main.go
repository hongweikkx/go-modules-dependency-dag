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

// en！ 深度是我向下调用了多深。
var dag map[string]map[string]int
func main() {
	// config
	pkgRootName := *flag.String("app", "beego", "package name")
	path := *flag.String("path", "/Users/gaohongwei/go/src/github.com/astaxie/beego", "app path")
	dotFilePath := *flag.String("dot-filepath", "./dag.dot", "dot out path")
	flag.Parse()
	path, _ = filepath.Abs(path)
	// check
	if !Exist(path) {
		fmt.Print("app path can not use")
		return
	}
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
			fmt.Printf("visited file or dir: %q %q\n",singlePath, info.Name())
			parseSingleFile(pkgRootName, path, singlePath)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", path, err)
	}
}

func parseSingleFile(pkgRootName, path, singlePath string) {
	f, err := os.Open(singlePath)
	if err != nil {
		panic(fmt.Sprintf("parse single file %s error:%s", singlePath, err.Error()))
	}
	r := bufio.NewReader(f)
	isImport := false
	packageName := packageNamePro(path, singlePath)
	for {
		line,_, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(fmt.Sprintf("read single file %s error:%s", singlePath, err.Error()))
		}
		if string(line) == "import (" {
			isImport = true
			continue
		}
		if isImport == true && string(line) == ")"{
			isImport = false
			break
		}
		fmt.Println("bbbbb:", string(line))
		if isImport == true && isSameRoot(string(line), path){
			fmt.Println("aaaaaa:", string(line))
			if _, ok := dag[packageName]; ok {
				dag[packageName][string(line)] = 1
			}else {
				tmp := make(map[string]int)
				tmp[string(line)] = 1
				dag[packageName] = tmp
			}
		}
	}
}

func writeToDotFile(dotFilePath string) {
}

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}


func packageNamePro(path, singlePath string) string {
	ret := ""
	pathSplit := strings.Split(path, "/")
	SinglePathSplit := strings.Split(singlePath, "/")
	for i:=len(pathSplit)-1; i < len(SinglePathSplit); i++{
		if i == len(SinglePathSplit) - 1 {
			ret += SinglePathSplit[i]
		}else {
			ret += SinglePathSplit[i] + "/"
		}
	}
	return ret
}

func isSameRoot(packageName, path string) bool {
	pathSplit := strings.Split(path, "/")
	singlePathSplit := strings.Split(packageName, "/")
	return pathSplit[len(pathSplit) - 1] == singlePathSplit[0]
}