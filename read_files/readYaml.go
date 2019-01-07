package read_files

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

//读取yaml文件
type TemplateYaml struct {
	StructName string  `json:"struct_name" yaml:"struct_name"`
	Fields     []Field `json:"fields" yaml:"fields"`
	HasTest    bool    `json:"has_test" yaml:"has_test"`
	MenuLevel  int     `json:"menu_level" yaml:"menu_level"`
}

type Field struct {
	FieldName  string `yaml:"field_name"`  //字段名
	FieldType  string `yaml:"field_type"`  //字段类型
	IsJson     bool   `yaml:"is_json"`     //是否为json数据
	FieldLabel string `yaml:"field_label"` //字段显示值
	IsHidden   bool   `yaml:"is_hidden"`   //是否隐藏
	JsonName   string `yaml:"json_name"`   //json格式名称
	//add、edit的表单类型  不填
	// 默认为text 普通文本框 可选值:
	// select 单选框
	// multiple_select: 多选框
	// date 日期选择框
	// datetime 日期+时间选择
	// email 邮箱控件
	// radio
	// checkbox
	// 目前只支持text 、 radio (bool类型值)、date、datetime
	DisplayType string `yaml:"display_type"`
	IndexWidth  int    `yaml:"index_width"`
	IsDisplay   bool   `yaml:"is_display"`
	Editable    bool   `yaml:"editable"`
	IsQuery     bool   `yaml:"is_query"`
	Comment     string `yaml:"comment"`
}

type TemplateStruct struct {
	ObjectName      string  `yaml:"object_name"`
	ObjectNames     string  `yaml:"object_names"`
	ObjectBsonName  string  `yaml:"object_bson_name"`
	ObjectBsonNames string  `yaml:"object_bson_names"`
	ObjectJsonName  string  `yaml:"object_json_name"`
	ObjectJsonNames string  `yaml:"object_json_names"`
	Label           string  `yaml:"label"`
	NotUpdate       bool    `yaml:"not_update"`
	NotDelete       bool    `yaml:"not_delete"`
	Comment         string  `yaml:"comment"`
	Fields          []Field `yaml:"fields"`
}

func ReadYaml(tempPath string) *TemplateStruct {
	bs, err := ioutil.ReadFile(tempPath)
	if err != nil {
		panic(err)
	}
	var t TemplateStruct
	err = yaml.Unmarshal(bs, &t)
	if err != nil {
		panic(err)
	}
	return &t
}
