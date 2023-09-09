package main

import "gorm.io/gorm"

// type MyStruct struct {
// 	gorm.Model
// 	// 注释1
// 	Field1 int

// 	Field2 string // 注释2
// }

type Class struct {
	gorm.Model
	SchoolLabel  int    `json:"schoolLabel" gorm:"size:30"`   // 学校名称
	ClassLabel   string `json:"classLabel" gorm:"size:30"`    // 班级名称
	Teachers     string `json:"teachers" gorm:"size:255"`     // 带队老师(可为多个老师用逗号隔开)
	ClassMonitor string `json:"classMonitor" gorm:"size:255"` // 班长
	Student      string `json:"student" gorm:"size:255"`      // 学生
}
