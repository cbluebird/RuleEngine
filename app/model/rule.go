package model

import "time"

type Rule struct {
	Message    string
	CreateTime time.Time
	ID         int `gorm:"primary_key;AUTO_INCREMENT"`
}

var GetModule map[int]string = map[int]string{1: "art", 2: "gpa", 3: "innovate", 4: "labour", 5: "moral", 6: "pe"}
