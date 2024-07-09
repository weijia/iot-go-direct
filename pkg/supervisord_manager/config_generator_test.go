package supervisord_manager

import (
	"testing"
)

func TestGenerateConfig(t *testing.T) {
	// 创建一个字符串到字符串的映射
	stringMap := map[string]string{
		"main": "main.exe",
	}
	GenerateConfig(stringMap)
}
