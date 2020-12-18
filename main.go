package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
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
	writeToDotFile(dotFilePath)
	fmt.Println("end...")
}

func parse(pkgRootName, path string) {
	err := filepath.Walk(path, func(singlePath string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", singlePath, err)
			return err
		}
		// skip vendor
		if info.IsDir() && info.Name() == "vendor" {
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
	f, err := os.Open(singlePath)
	if err != nil {
		panic(fmt.Sprintf("parse single file %s error:%s", singlePath, err.Error()))
	}
	r := bufio.NewReader(f)
	isImport := false
	packageName := ""
	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(fmt.Sprintf("read single file %s error:%s", singlePath, err.Error()))
		}
		if strings.HasPrefix(string(line), "package") {
			packageNames := strings.Split(string(line), " ")
			packageName += filepath.Join(packageName, packageNames[len(packageNames)-1])
		}
		if string(line) == "import (" {
			isImport = true
			continue
		}
		if isImport == true && string(line) == ")" {
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
			} else {
				tmp := make(map[string]int)
				tmp[importName] = 1
				dag[packageName] = tmp
			}
		}
	}
}

func writeToDotFile(dotFilePath string) {
	fd, _ := os.OpenFile(dotFilePath, os.O_RDWR|os.O_CREATE, 0644)
	indegree := make(map[string]int)
	defer fd.Close()
	fd.Write([]byte("digraph G {\n"))
	for k, vm := range dag {
		for v := range vm {
			indegree[v]++
			fd.Write([]byte(fmt.Sprintf("\t\"%s\" -> \"%s\"\n", k, v)))
		}
	}
	colors := colorUseIndegree(indegree)
	//"A" [shape=circle, style=filled, fillcolor=red]
	for i := 0; i < len(colors); i++ {
		fd.Write([]byte(fmt.Sprintf("\t\"%s\" [fillcolor=\"%s\", style=filled]\n", colors[i].Label, colors[i].Color)))
	}

	fd.Write([]byte("}\n"))
}

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func importPkgName(line, pkgRootName string) string {
	if index := strings.Index(line, pkgRootName); index != -1 {
		name := strings.Trim(line[index:], "\"")
		names := strings.Split(name, "/")
		return names[len(names)-1]
	}
	return ""

}

type colorNode struct {
	Label    string
	Color    string
	InDegree int
}

type colorNodes []colorNode

func colorUseIndegree(indegree map[string]int) colorNodes {
	var nodes colorNodes
	for k, v := range indegree {
		nodes = append(nodes, colorNode{Label: k, InDegree: v})
	}
	sort.Sort(nodes)
	g := 255
	for i := 0; i < len(nodes); i++ {
		g = g - 255/len(nodes)
		nodes[i].Color = rgb2hex(255, int64(g), 0)
	}
	return nodes
}

func (nodes colorNodes) Len() int {
	return len(nodes)
}

func (nodes colorNodes) Swap(i, j int) {
	nodes[i], nodes[j] = nodes[j], nodes[i]
}

func (nodes colorNodes) Less(i, j int) bool {
	return nodes[i].InDegree < nodes[j].InDegree
}

// rgb -> hex
func rgb2hex(r, g, b int64) string {
	r16 := t2x(r)
	g16 := t2x(g)
	b16 := t2x(b)
	return "#" + r16 + g16 + b16
}

func t2x(t int64) string {
	result := strconv.FormatInt(t, 16)
	if len(result) == 1 {
		result = "0" + result
	}
	return result
}
