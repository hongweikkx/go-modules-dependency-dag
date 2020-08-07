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
	packageName := *flag.String("app", "beego", "package name")
	path := *flag.String("path", "/Users/gaohongwei/go/src/github.com/astaxie/beego", "app path")
	depth := *flag.Int("dep", 10, "modules depth")
	dotFilePath := *flag.String("dot-filepath", "./dag.dot", "dot out path")
	flag.Parse()
	path, _ = filepath.Abs(path)
	// check
	if !Exist(path) {
		fmt.Print("app path can not use")
		return
	}
	if depth < 1 {
		fmt.Printf("depth is not valid. must bigger than 0")
		return
	}
	fmt.Printf("path: %s, depth: %d, dot filepath: %s \n", path, depth, dotFilePath)
	fmt.Println("start...")
	dag = make(map[string]map[string]int)
	parse(path, depth)
	for k, vm := range dag {
		for vk, _ := range vm {
			fmt.Println(k, vk)
		}
	}
	writeToDotFile(dotFilePath)
	fmt.Println("end...")
}


func parse(path string, depth int) {
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
		// dep filter
		if info.IsDir() &&  pathDepth(singlePath) - pathDepth(path) >= depth{
			//fmt.Printf("skipping a dir with depth limit: %+v \n", singlePath)
			return filepath.SkipDir
		}
		// select go file
		if !info.IsDir() && strings.Contains(info.Name(), ".go") {
			fmt.Printf("visited file or dir: %q %q\n",singlePath, info.Name())
			parseSingleFile(singlePath)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", path, err)
	}
}

func parseSingleFile(singlePath string) {
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
			packList := strings.Split(string(line), " ")
			packageName = packList[len(packList) -1]
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

// return the path depth
func pathDepth(path string) int{
	depth := 0
	for i := 0; i < len(path); i++ {
		if os.IsPathSeparator(path[i]) {
			depth++
		}
	}
	return depth
}

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
