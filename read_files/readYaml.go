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
	FieldName   string `yaml:"field_name"`
	FieldType   string `yaml:"field_type"`
	IsJson      bool   `yaml:"is_json"`
	FieldLabel  string `yaml:"field_label"`
	IsHidden    bool   `yaml:"is_hidden"`
	JsonName    string `yaml:"json_name"`
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
