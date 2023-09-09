package main

import "gorm.io/gorm"

type MyStruct struct {
	gorm.Model
	// 注释1
	Field1 int

	Field2 string // 注释2
}
