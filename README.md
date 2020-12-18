[![Go Report Card](https://goreportcard.com/badge/github.com/hongweikkx/go-modules-dependency-dag)](https://goreportcard.com/report/github.com/hongweikkx/go-modules-dependency-dag)

# go-modules-dependency-dag
Analyze the dependencies of go modules (except third-party packages) under the specified path. Use colors to distinguish the frequency of being imported by other modules.From yellow to red, Frequency increase

## Install
```
go get -u github.com/hongweikkx/go-modules-dependency-dag
```

## Run
```
go-modules-dependency-dag --app="app name" --path="app path" --dot-filepath="dot filepath"
```

example:

```
go-modules-dependency-dag --app=github.com/astaxie/beego --path=/Users/hongweigaokkx/go/src/github.com/astaxie/beego --dot_file_path=dag.dot
```


## Visualization
use graphviz to visualize dot file.

```
dot -T png dag.dot > dag.png
dot -T svg dag.dot > dag.svg
```

example:
![image](https://github.com/hongweikkx/go-modules-dependency-dag/blob/master/example/dag.png)

## refer
* [https://github.com/legendtkl/godag](https://github.com/legendtkl/godag)
* [Graphviz](http://www.graphviz.org/)
