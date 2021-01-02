package config

type LoggerConfig struct {
	File *LoggerConfig_File `json:"file"`
}

type LoggerConfig_File struct {
	Prefix   string `json:"prefix"`
	Suffix   string `json:"suffix"`
	FileName string `json:"filename"`
	Daily    bool   `json:"daily"`
	Level    int    `json:"level"`
}
