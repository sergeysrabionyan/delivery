package excel_builder

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

const defaultSheet = "TDSheet"

type Builder struct {
	CheckFilePath    string
	DecipherFilePath string
}

type SubBuilder struct {
	file  *excelize.File
	sheet string
}

func New() *Builder {
	return &Builder{}
}

func (b *Builder) Validate() bool {
	if b.CheckFilePath == "" {
		return false
	}
	if b.DecipherFilePath == "" {
		return false
	}
	return true
}

func (b *Builder) GetCheckBuilder() (*SubBuilder, *excelize.File) {
	f, err := excelize.OpenFile(b.CheckFilePath)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	return getSubBuilder(f), f
}

func (b *Builder) GetDecipherBuilder() (*SubBuilder, *excelize.File) {
	f, err := excelize.OpenFile(b.DecipherFilePath)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	return getSubBuilder(f), f
}

func getSubBuilder(file *excelize.File) *SubBuilder {
	return &SubBuilder{file: file, sheet: defaultSheet}
}

func (b *SubBuilder) ChangeCell(cell string, value interface{}) {
	err := b.file.SetCellValue(b.sheet, cell, value)
	if err != nil {
		panic(err)
	}
}

func (b *SubBuilder) ClearCell(cell string) {
	err := b.file.SetCellValue(b.sheet, cell, "")
	if err != nil {
		panic(err)
	}
	err = b.file.SetCellFormula(b.sheet, cell, "")
	if err != nil {
		panic(err)
	}
}

func (b *SubBuilder) GetCellValue(cell string) string {
	value, err := b.file.GetCellValue(b.sheet, cell)
	if err != nil {
		return ""
	}
	return value
}

func (b *SubBuilder) Save() error {
	return b.file.Save()
}
