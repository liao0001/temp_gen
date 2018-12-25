package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"template/read_files"
	"text/template"
)

var mainRoot string
var packageName string
var selfPath string

func init() {
	flag.StringVar(&mainRoot, "p", "", "revel项目根目录的绝对路径")
	flag.Parse()
	mainRoot = `E:\workSpace\go\dev\src\test_template`

	//checkRevelProject(mainRoot)
	getTempPath()
}

func main() {
	//checkRevelProject(mainRoot)
	//test()

	testModel()
}

//测试模板
func testModel() {
	obj := read_files.ReadYaml("./roles.yaml")

	filePath := filepath.Join(selfPath, obj.ObjectJsonNames+".go")
	os.Remove(filePath)

	f, err := os.Create(filePath)
	if err != nil {
		panic("创建文件失败," + err.Error())
	}

	res := map[string]interface{}{
		"obj": obj,
	}

	tempPath := filepath.Join(selfPath, "temps", "db.temp")
	te := getTemplate(tempPath, obj.ObjectBsonName)

	err = te.Execute(f, res)
	if err != nil {
		panic("转化模板失败," + err.Error())
	}
}

//获取模板
func getTemplate(p string, name string) *template.Template {
	bs, err := ioutil.ReadFile(p)
	if err != nil {
		panic("读取文件内容失败," + err.Error())
	}

	t, err := template.New(name).Parse(string(bs))
	if err != nil {
		panic("解析模板失败," + err.Error())
	}
	return t
}

//检测是否为revel项目
func checkRevelProject(mainRoot string) {
	//检测是否为revel项目
	fileInfo, err := os.Stat(mainRoot)
	if err != nil || !fileInfo.IsDir() {
		panic("无效路径:" + mainRoot)
	}

	//检查必要目录或者文件
	dirs := []string{
		filepath.Join(mainRoot, "app", "controllers"),
		filepath.Join(mainRoot, "app", "views"),
		filepath.Join(mainRoot, "app", "init.go"),
		filepath.Join(mainRoot, "conf", "app.conf"),
		filepath.Join(mainRoot, "public"),
		filepath.Join(mainRoot, "templates"),
	}
	for _, fn := range dirs {
		_, err = os.Stat(fn)
		if err != nil {
			panic("项目缺少文件或者目录: " + fn)
		}
	}
}

func getTempPath() {
	p, err := os.Getwd()
	if err != nil {
		panic("获取当前目录信息失败，" + err.Error())
	}
	selfPath = p
}

//测试模板
func test() {
	tempPath := filepath.Join(selfPath, "temps", "models", "test.temp")
	bs, err := ioutil.ReadFile(tempPath)
	if err != nil {
		panic("读取文件内容失败," + err.Error())
	}
	t, err := template.New("sss").Parse(string(bs))
	if err != nil {
		panic("解析模板失败," + err.Error())
	}

	filePath := filepath.Join(selfPath, "roles.go")
	os.Remove(filePath)

	f, err := os.Create(filePath)
	if err != nil {
		panic("创建文件失败," + err.Error())
	}

	arr := []struct {
		Name string
		Age  int
	}{
		{Name: "测试名称", Age: 1},
		{Name: "测试名称2", Age: 2},
		{Name: "测试名称3", Age: 3},
	}
	res := map[string]interface{}{
		"arr": arr,
	}

	err = t.Execute(f, res)
	if err != nil {
		panic("转化模板失败," + err.Error())
	}
}
