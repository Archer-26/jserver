package main

import (
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
)

var (
	lineNumber     = 2                     // 每个工作表需要读取的行数
	structBegin    = "type _%v struct {\n" // 结构体开始
	structValue    = "    %v %v"           // 结构体的内容
	structRemarks  = "	 // %v"             // 结构体备注
	structValueEnd = "\n"                  // 结构体内容结束
	structEnd      = "}\n"                 // 结构体结束
	headerFromat   = "package %v\n\r"      // 文件头
)

type Generate struct {
	SaveGoPath   string // 生成go文件的保存路径
	SaveJsonPath string // 生成json文件的保存路径
	TplPath      string // 模板文件路径
	ReadPath     string // excel表路径
	PackageName  string // 包名
}

// 读取excel
func (this *Generate) ReadExcel() {
	files, err := ioutil.ReadDir(this.ReadPath)
	if err != nil {
		panic(fmt.Errorf("excel文件路径读取失败 此路径无效:%v error:%v", this.ReadPath, err))
	}
	var wg sync.WaitGroup
	for i, file := range files {
		if hasChinese(file.Name()) {
			continue
		}
		if strings.Contains(file.Name(), "~$") {
			continue
		}
		if !strings.Contains(file.Name(), ".xlsx") {
			continue
		}

		wg.Add(1)
		go func(j int) {
			defer wg.Done()

			dir := path.Join(this.ReadPath, files[j].Name())
			f, err := xlsx.OpenFile(dir)
			if err != nil {
				panic(fmt.Errorf("excel文件读取失败 无效文件:%v error:%v", dir, err))
			}

			// 遍历工作表
			for _, sheet := range f.Sheets {
				if err := this.BuildTypeStruct(sheet, files[j].Name()); err != nil {
					panic(err)
				}
				if err := this.BuildJsonStruct(sheet, files[j].Name()); err != nil {
					panic(err)
				}
			}
		}(i)
	}
	wg.Wait()

	// 导出依赖的字段解析函数
	tpath := path.Join(this.TplPath, "convert.go.tpl")
	tmpl, err := template.ParseFiles(tpath)
	if err != nil {
		panic(fmt.Errorf("模板文件读取失败，无效路径:%v err:%v", tpath, err))
	}

	newFile := path.Join(this.SaveGoPath, "convert.funtype.go")
	file, err := os.OpenFile(newFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer file.Close()
	if err != nil {
		panic(fmt.Errorf("配置文件创建失败 无效的文件或路径:%v err:%v", newFile, err))
	}

	err = tmpl.Execute(file, struct {
		Packagename string
	}{this.PackageName})

	if err != nil {
		panic(fmt.Errorf("模板执行输出失败 err:%v", err))
	}
	exec.Command("go", "fmt", this.SaveGoPath).Run()
}

// 构建类型结构
func (this *Generate) BuildTypeStruct(sheet *xlsx.Sheet, fileName string) error {
	sheetData := make([][]string, 0)
	// 判断表格中内容的行数是否小于需要读取的行数
	if sheet.MaxRow < lineNumber {
		return fmt.Errorf("ReadExcel sheet.MaxRow:%d < lineNumber:%d file:%v", sheet.MaxRow, lineNumber, fileName)
	}
	// 遍历列
	for i := 0; i < sheet.MaxCol; i++ {
		// 没有字段，不解析
		if strings.TrimSpace(sheet.Cell(FIELDNAME, i).Value) == "" {
			continue
		}
		cellData := make([]string, 0)
		// 遍历行
		for j := 0; j < lineNumber; j++ {
			cellData = append(cellData, strings.TrimSpace(sheet.Cell(j, i).Value))
		}
		sheetData = append(sheetData, cellData)
	}
	structType, err := parsing(sheetData, sheet.Name)
	if err != nil {
		return fmt.Errorf("fileName:\"%v\" is err:%v", fileName, err)
	}

	if structType == "" {
		return fmt.Errorf("ReadExcel this.data is nil")
	}
	structType += "\n"

	fieldName := strings.TrimSpace(sheet.Cell(FIELDNAME, 0).Value)
	sp := strings.Split(fieldName, "_")
	if len(sp) < 1 {
		return fmt.Errorf("字段名格式错误:%v %v", fileName, sheet.Name)
	}
	var keyType string
	fieldType := sp[0]
	if k, ex := TypeIndex[fieldType]; !ex {
		return fmt.Errorf("主键类型找不到:%v %v %v", fieldType, fileName, sheet.Name)
	} else {
		keyType = k
	}
	if keyType != "string" && keyType != "int64" {
		return fmt.Errorf("主键类型错误:%v %v %v", fieldType, fileName, sheet.Name)
	}

	err = this.writeGolangFile(structType, sheet.Name, keyType, fieldName, fileName)
	if err != nil {
		return err
	}
	return nil
}

// 构建json结构
func (this *Generate) BuildJsonStruct(sheet *xlsx.Sheet, fileName string) error {
	array := []string{}
	// 判断表格中内容的行数是否小于需要读取的行数
	dataLen := sheet.MaxRow - lineNumber
	if dataLen < 0 {
		return fmt.Errorf("ReadExcel dataLen < 0 dataLen:%v MaxRow:%v lineNumber:%v ", dataLen, sheet.MaxRow, lineNumber)
	}
	// 遍历列
	var err error
	for i := lineNumber; i <= sheet.MaxRow; i++ {
		// 遍历行
		cell := sheet.Cell(i, 0)
		if cell == nil {
			return fmt.Errorf("fileName:%v cell == nil ", fileName)
		}

		// 单独处理key
		key := strings.TrimSpace(cell.Value)
		if key == "" {
			break
		}
		m := map[string]interface{}{}
		for j := 0; j < sheet.MaxCol; j++ {
			if strings.TrimSpace(sheet.Cell(1, j).Value) == "" {
				continue
			}
			fieldName := strings.TrimSpace(sheet.Cell(FIELDNAME, j).Value)
			sp := strings.Split(fieldName, "_")
			if len(sp) < 1 {
				return fmt.Errorf("字段名格式错误:%v", fieldName)
			}
			fieldType := sp[0]
			if _, ex := TypeIndex[fieldType]; !ex {
				continue
			}

			if m[fieldName], err = TypeConvert[fieldType](strings.TrimSpace(sheet.Cell(i, j).Value)); err != nil {
				return fmt.Errorf("TypeConvert error=%v i=%v  j=%v name=%v value=%v file=%v", err, i, j, fieldName, sheet.Cell(i, j).Value,fileName)
			}
		}
		if data, err := json.Marshal(m); err == nil {
			array = append(array, string(data))
		} else {
			return fmt.Errorf("json.Marshal(array) error:%v ", err)
		}
	}
	if err := this.writeJsonFile("[\n    "+strings.Join(array, ",\n    ")+"\n]", sheet.Name); err != nil {
		return err
	}
	return nil
}
