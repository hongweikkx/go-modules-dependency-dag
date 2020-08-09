# go-modules-dependency-dag
分析指定路径下的go模块(除第三方包之外)依赖。<br>
并用颜色区分被其他模块import的频次。从黄色到红色，频次递增。 

## Install
go get -u github.com/hongweikkx/go-modules-dependency-dag

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
