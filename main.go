package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"temp_gen/read_files"
	"text/template"
)

var mainRoot string
var packageName string
var selfPath string

type TargetPath struct {
	ModelsPath      string `yaml:"models_path"`
	ServicesPath    string `yaml:"services_path"`
	ControllersPath string `yaml:"controllers_path"`
	WebsDir         string `yaml:"webs_dir"`
	JsDir           string `yaml:"js_dir"`
	RouterFilePath  string `yaml:"router_file_path"`
	Types           string
	Forced          bool
	SourceFile      string
}

var globalPathConfig TargetPath

func init() {
	bs, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(bs, &globalPathConfig)
	if err != nil {
		panic(err)
	}

	var h string
	flag.StringVar(&h, "h", "", "-h 查看参数含义")

	var t string
	flag.StringVar(&t, "t", "all", "types 表示需要生成的文件类型，可选项为: all(默认),js,webs,db,service,controller")

	var s string
	flag.StringVar(&s, "s", "", "配置文件地址（yaml格式） 不能为空")

	var f bool
	flag.BoolVar(&f, "f", false, "是否强制替换; false(默认): 如果文件已存在，则会给出提示并跳过;true: 会强制覆盖原文件(慎用)")
	flag.Parse()
	globalPathConfig.Forced = f
	globalPathConfig.SourceFile = s
	globalPathConfig.Types = t
	if len(strings.TrimSpace(s)) == 0 {
		panic("指定配置文件不能为空,-h 查看帮助")
	}
}

func main() {
	//checkRevelProject(mainRoot)
	//test()

	generateTemps(globalPathConfig.SourceFile)
}

//测试模板
func generateTemps(targetFile string) {
	var err error
	obj := read_files.ReadYaml(targetFile)
	if globalPathConfig.Types == "all" || globalPathConfig.Types == "db" {
		err = renderDb(obj, globalPathConfig.Forced)
		if err != nil {
			fmt.Printf("生成db文件失败:" + err.Error())
		}
	}
	if globalPathConfig.Types == "all" || globalPathConfig.Types == "service" {
		err = renderService(obj, globalPathConfig.Forced)
		if err != nil {
			fmt.Printf("生成service文件失败:" + err.Error())
		}
	}
	if globalPathConfig.Types == "all" || globalPathConfig.Types == "controller" {
		err = renderController(obj, globalPathConfig.Forced)
		if err != nil {
			fmt.Printf("生成controller文件失败:" + err.Error())
		}
	}

	if globalPathConfig.Types == "all" || globalPathConfig.Types == "js" {
		err = renderJs(obj, globalPathConfig.Forced)
		if err != nil {
			fmt.Printf("生成js文件失败:" + err.Error())
		}
	}

	if globalPathConfig.Types == "all" || globalPathConfig.Types == "webs" {
		err = renderWebs(obj, globalPathConfig.Forced)
		if err != nil {
			fmt.Printf("生成html文件失败:" + err.Error())
		}
	}

	updateRouterFile(obj)
}

func renderDb(templ *read_files.TemplateStruct, forced bool) error {
	filePath := filepath.Join(globalPathConfig.ModelsPath, templ.ObjectJsonNames+".go")
	tempName := filepath.Join("apps", "db.temp")
	//生成对应的repo 和 sync
	updateDBFile(templ)
	return renderFile(templ, filePath, tempName, forced)
}

func renderService(templ *read_files.TemplateStruct, forced bool) error {
	filePath := filepath.Join(globalPathConfig.ServicesPath, templ.ObjectJsonName+"_service.go")
	tempName := filepath.Join("apps", "service.temp")
	updateServiceFile(templ)
	return renderFile(templ, filePath, tempName, forced)
}

func renderController(templ *read_files.TemplateStruct, forced bool) error {
	filePath := filepath.Join(globalPathConfig.ControllersPath, templ.ObjectJsonNames+".go")
	tempName := filepath.Join("apps", "controller.temp")
	return renderFile(templ, filePath, tempName, forced)
}

func renderJs(templ *read_files.TemplateStruct, forced bool) error {
	//创建目录
	dirPath := filepath.Join(globalPathConfig.JsDir, templ.ObjectJsonNames)
	os.MkdirAll(dirPath, 0777)

	var err error
	indexJs := filepath.Join(dirPath, "index.js")
	indexTempName := filepath.Join("js", "index.js.temp")
	err = renderFile(templ, indexJs, indexTempName, forced)
	if err != nil {
		fmt.Println("生成index.js文件失败," + err.Error())
	}

	addJs := filepath.Join(dirPath, "add.js")
	addTempName := filepath.Join("js", "add.js.temp")
	err = renderFile(templ, addJs, addTempName, forced)
	if err != nil {
		fmt.Println("生成add.js文件失败," + err.Error())
	}

	editJs := filepath.Join(dirPath, "edit.js")
	editTempName := filepath.Join("js", "edit.js.temp")
	err = renderFile(templ, editJs, editTempName, forced)
	if err != nil {
		fmt.Println("生成edit.js文件失败," + err.Error())
	}
	return nil
}

func renderWebs(templ *read_files.TemplateStruct, forced bool) error {
	//创建目录
	dirPath := filepath.Join(globalPathConfig.WebsDir, templ.ObjectJsonNames)
	os.MkdirAll(dirPath, 0777)

	var err error
	indexHtml := filepath.Join(dirPath, "index.html")
	indexTempName := filepath.Join("htmls", "index.html.temp")
	err = renderFile(templ, indexHtml, indexTempName, forced)
	if err != nil {
		fmt.Println("生成index.html文件失败," + err.Error())
	}

	addHtml := filepath.Join(dirPath, "add.html")
	addTempName := filepath.Join("htmls", "add.html.temp")
	err = renderFile(templ, addHtml, addTempName, forced)
	if err != nil {
		fmt.Println("生成add.html文件失败," + err.Error())
	}

	editHtml := filepath.Join(dirPath, "edit.html")
	editTempName := filepath.Join("htmls", "edit.html.temp")
	err = renderFile(templ, editHtml, editTempName, forced)
	if err != nil {
		fmt.Println("生成edit.html文件失败," + err.Error())
	}
	return nil
}

func renderFile(templ *read_files.TemplateStruct, filePath string, tempName string, forced bool) error {
	if forced {
		os.Remove(filePath)
	} else if _, err := os.Stat(filePath); err == nil {
		fmt.Printf("文件已存在:%s\n", filePath)
		os.Exit(-1)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return errors.New("创建文件失败," + err.Error())
	}

	res := map[string]interface{}{
		"obj": templ,
	}

	tempPath := filepath.Join(selfPath, "temps", tempName)
	te := getTemplate(tempPath, templ.ObjectBsonName)

	err = te.Execute(f, res)
	if err != nil {
		return errors.New("转化模板失败," + err.Error())
	}
	return nil
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

//单独处理db.go 文件
func updateDBFile(templ *read_files.TemplateStruct) {
	filePath := filepath.Join(globalPathConfig.ModelsPath, "db.go")
	repoSignStr := "////// repo"
	syncTableStr := "////// syncTable"

	repoNewStr := fmt.Sprintf("func (db *DB) %s() *%sRepo {\n"+
		"\treturn &%sRepo{orm.New(func() interface{} {\n"+
		"\t\treturn &%s{}\n"+
		"\t})(db.Engine).WithSession(db.session)}\n"+
		"}\n\n", templ.ObjectNames, templ.ObjectName, templ.ObjectName, templ.ObjectName)

	syncNewStr := fmt.Sprintf("\t\t&%s{},\n", templ.ObjectName)

	bs, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("读取db.go失败:" + err.Error())
		return
	}

	needUpdate := false
	if bytes.Index(bs, []byte(repoNewStr)) == -1 {
		bs = bytes.Replace(bs, []byte(repoSignStr), []byte(repoNewStr+repoSignStr), -1)
		needUpdate = true
	}

	if bytes.Index(bs, []byte(syncNewStr)) == -1 {
		bs = bytes.Replace(bs, []byte(syncTableStr), []byte(syncNewStr+syncTableStr), -1)
		needUpdate = true
	}

	if needUpdate {
		err = ioutil.WriteFile(filePath, bs, 0666)
		if err != nil {
			fmt.Println("重写db.go失败:" + err.Error())
			return
		}
	}

	updateIndexFields(templ)
}

func updateServiceFile(templ *read_files.TemplateStruct) {
	filePath := filepath.Join(globalPathConfig.ServicesPath, "base_service.go")
	signStr := "////// service defined"
	setStr := "////// service set"

	signNewStr := fmt.Sprintf("\n\t%sService %sService\n", templ.ObjectName, templ.ObjectName)

	setNewStr := fmt.Sprintf("\t%sService{baseService},\n", templ.ObjectName)

	bs, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("读取db.go失败:" + err.Error())
		return
	}

	needUpdate := false
	if bytes.Index(bs, []byte(signNewStr)) == -1 {
		bs = bytes.Replace(bs, []byte(signStr), []byte(signNewStr+signStr), -1)
		needUpdate = true
	}

	if bytes.Index(bs, []byte(setNewStr)) == -1 {
		bs = bytes.Replace(bs, []byte(setStr), []byte(setNewStr+setStr), -1)
		needUpdate = true
	}

	if needUpdate {
		err = ioutil.WriteFile(filePath, bs, 0666)
		if err != nil {
			fmt.Println("重写base_service.go失败:" + err.Error())
			return
		}
	}
}

func updateRouterFile(templ *read_files.TemplateStruct) {
	filePath := globalPathConfig.RouterFilePath
	signStr := "########## template sign"

	signNewStr := fmt.Sprintf("\n"+
		"# %s\n"+
		"GET     /api/%s/list                       %s.List\n"+
		"GET     /api/%s/get                        %s.GetById\n"+
		"POST    /api/%s/add                        %s.Create\n"+
		"POST    /api/%s/update                     %s.Update\n"+
		"POST    /api/%s/delete                     %s.Delete\n"+
		"POST    /api/%s/deletebyids                %s.DeleteByIds\n\n", templ.Label, templ.ObjectJsonNames, templ.ObjectNames, templ.ObjectJsonNames, templ.ObjectNames,
		templ.ObjectJsonNames, templ.ObjectNames, templ.ObjectJsonNames, templ.ObjectNames,
		templ.ObjectJsonNames, templ.ObjectNames, templ.ObjectJsonNames, templ.ObjectNames)

	bs, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("读取router失败:" + err.Error())
		return
	}

	needUpdate := false
	if bytes.Index(bs, []byte(signNewStr)) == -1 {
		bs = bytes.Replace(bs, []byte(signStr), []byte(signNewStr+signStr), -1)
		needUpdate = true
	}

	if needUpdate {
		err = ioutil.WriteFile(filePath, bs, 0666)
		if err != nil {
			fmt.Println("重写router失败:" + err.Error())
			return
		}
	}
}

func updateIndexFields(templ *read_files.TemplateStruct) {
	filePath := filepath.Join(globalPathConfig.ControllersPath, "../libs", "init_data.go")
	signStr := "////// index fields"
	signNewStr := fmt.Sprintf("\terr = db.IndexFieldSettings().InitFieldInfoForTable(db.%s().Name(), \"%s\", db.%s().GetFields())\n"+
		"\tif err != nil {\n"+
		"\t\treturn err\n"+
		"\t}\n\n", templ.ObjectNames, templ.Label, templ.ObjectNames)
	bs, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("读取router失败:" + err.Error())
		return
	}

	needUpdate := false
	if bytes.Index(bs, []byte(signNewStr)) == -1 {
		bs = bytes.Replace(bs, []byte(signStr), []byte(signNewStr+signStr), -1)
		needUpdate = true
	}

	if needUpdate {
		err = ioutil.WriteFile(filePath, bs, 0666)
		if err != nil {
			fmt.Println("重写router失败:" + err.Error())
			return
		}
	}
}
