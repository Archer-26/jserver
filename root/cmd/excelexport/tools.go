package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"unicode"
)

const COMMENT = 0   // 字段注释
const FIELDNAME = 1 // 字段名

// 解析配置表格式
func parsing(data [][]string, structName string) (result string, err error) {
	stName := firstRuneToUpper(structName)
	structData := fmt.Sprintf(structBegin, stName)
	fieldGet := ""
	for _, value := range data {
		if len(value) != lineNumber {
			return "", fmt.Errorf("parsing sheetName:%v col's len:%d is err", value, len(value))
		}

		comment := value[COMMENT]
		fieldName := value[FIELDNAME]
		sp := strings.Split(fieldName, "_")
		if len(sp) < 1 {
			return "", fmt.Errorf("字段名格式错误:%v", fieldName)
		}
		fieldType := sp[0]
		_, exist := TypeIndex[fieldType]
		if !exist {
			continue // 所有找不到的类型，全部忽略
			//return "", fmt.Errorf("parsing TypeIndex err:%v structName:%v", fieldType, structName)
		}

		goType := TypeIndex[fieldType]
		field := firstRuneToUpper(fieldName)
		structData += fmt.Sprintf(structValue, field, goType)
		fieldGet += genGet(stName, field, fieldType, comment)
		if comment != "" {
			structData += fmt.Sprintf(structRemarks, comment)
		}
		structData += fmt.Sprintf(structValueEnd)
	}
	structData += structEnd
	structData += fieldGet
	return structData, nil
}

func genGet(stName, field, fieldType, comment string) string {
	filterTypeField := strings.Join(strings.Split(field, "_")[1:], "_")
	filterTypeField = firstRuneToUpper(filterTypeField)
	switch fieldType {
	case INT, INT64, STR, FLOAT:
		return fmt.Sprintf("func (c *%v)%v() %v {return c.data.%v}//%v\n", stName, filterTypeField, ValueIndex[fieldType], field, comment)
	case ARRAYINT, ARRAYSTR, ARRAYFLOAT:
		len := fmt.Sprintf("//%v\nfunc (c *%v)Len_%v() int {return c.data.%v.Len()}\n", comment, stName, filterTypeField, field)
		get := fmt.Sprintf("func (c *%v)Get_%v(key %v) %v {return c.data.%v.Get(key)}\n", stName, filterTypeField, KeyIndex[fieldType], ValueIndex[fieldType], field)
		rg := fmt.Sprintf("func (c *%v)Range_%v(fn func(%v, %v) (stop bool)) {c.data.%v.Range(fn)}\n", stName, filterTypeField, KeyIndex[fieldType], ValueIndex[fieldType], field)
		cp := fmt.Sprintf("func (c *%v)Copy_%v() %v {return c.data.%v.Copy()}\n", stName, filterTypeField, TypeIndex[fieldType], field)
		return len + get + rg + cp
	case INT2INT, INT2STR, STR2INT, STR2STR:
		len := fmt.Sprintf("//%v\nfunc (c *%v)Len_%v() int {return c.data.%v.Len()}\n", comment, stName, filterTypeField, field)
		get := fmt.Sprintf("func (c *%v)Get_%v(key %v) (%v, bool) {return c.data.%v.Get(key)}\n", stName, filterTypeField, KeyIndex[fieldType], ValueIndex[fieldType], field)
		rg := fmt.Sprintf("func (c *%v)Range_%v(fn func(%v, %v) (stop bool)) {c.data.%v.Range(fn)}\n", stName, filterTypeField, KeyIndex[fieldType], ValueIndex[fieldType], field)
		cp := fmt.Sprintf("func (c *%v)Copy_%v() %v {return c.data.%v.Copy()}\n", stName, filterTypeField, TypeIndex[fieldType], field)
		return len + get + rg + cp
	}
	return ""
}

// 拼装好的struct写入新的.conf.go文件
func (this *Generate) writeGolangFile(struType, sheetName, keyType, key, xlsxname string) error {
	type Portion struct {
		FileHeaderComment string
		Packagename       string
		TypeName          string
		SheetName         string
		MapType           string
		ArrType           string
		Key               string
		KeyType           string
		StruType          string
	}
	format := Portion{}
	format.FileHeaderComment = xlsxname
	format.Packagename = fmt.Sprintf(headerFromat, *goPackageName)
	format.TypeName = firstRuneToUpper(sheetName)
	format.SheetName = sheetName
	format.StruType = struType
	format.KeyType = keyType
	format.Key = key

	// 包名+头注释+整体
	tpath := path.Join(this.TplPath, "config.go.tpl")
	tmpl, err := template.ParseFiles(tpath)
	if err != nil {
		return fmt.Errorf("模板文件读取失败，无效路径:%v err:%v", tpath, err)
	}

	newFile := path.Join(this.SaveGoPath, sheetName+".conf.go")
	fw, err := os.OpenFile(newFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer func() {
		fw.Close()
	}()
	if err != nil {
		return fmt.Errorf("配置文件创建失败 无效的文件或路径:%v err:%v", newFile, err)
	}

	err = tmpl.Execute(fw, format)
	if err != nil {
		return fmt.Errorf("模板执行输出失败:%v err:%v", err)
	}
	fmt.Printf("%-13v %v \n", "build golang", newFile)
	return nil
}

// 拼装好json写入新的.json 文件
func (this *Generate) writeJsonFile(data, filename string) error {
	newFile := path.Join(this.SaveJsonPath, filename+".json")
	if err := ioutil.WriteFile(newFile, []byte(data), 0644); err != nil {
		return fmt.Errorf("writeJsonFile Write is err:%v", err)
	}
	fmt.Printf("%-13v %v \n", "build json", newFile)
	return nil
}

// 字符串首字母转换成大写
func firstRuneToUpper(str string) string {
	data := []byte(str)
	for k, v := range data {
		if k == 0 {
			first := []byte(strings.ToUpper(string(v)))
			newData := data[1:]
			data = append(first, newData...)
			break
		}
	}
	return string(data[:])
}

// 判断是否存在汉字或者是否为全局的工作表
func hasChinese(r string) bool {
	for _, v := range []rune(r) {
		if unicode.Is(unicode.Han, v) {
			return true
		}
	}
	return false
}
