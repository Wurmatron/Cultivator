package model

type Configuration struct {
	Blockchain string
	Type       string
	Key        string
	Value      string
	LastUpdate int64
}

type Metrics struct {
	Blockchain string
	Style      string
	EntryType  string
	Type       string
	Value      string
	Timestamp  string
}
