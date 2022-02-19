package model

import "gorm.io/gorm"

type Configuration struct {
	gorm.Model
	Blockchain string
	Type       string
	Key        string
	Value      string
	LastUpdate int64
}

type Metrics struct {
	gorm.Model
	Blockchain string
	Style      string
	EntryType  string
	Type       string
	Value      string
	Timestamp  string
}
