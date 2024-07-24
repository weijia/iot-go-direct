package app_manager

import (
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 定义一个结构体实现 fs.DirEntry 接口
type fakeDirEntry struct {
    name string
}

// 实现 Name 方法
func (d fakeDirEntry) Name() string {
    return d.name
}

// 实现 IsDir 方法
func (d fakeDirEntry) IsDir() bool {
    return false
}

// 实现 Type 方法
func (d fakeDirEntry) Type() fs.FileMode {
    return 0
}

// 实现 Info 方法
func (d fakeDirEntry) Info() (fs.FileInfo, error) {
    return nil, nil
}

// 自定义的 ReadDir 函数，从字符串列表中获取文件名
func fakeReadDir(name string) ([]fs.DirEntry, error) {
    var entries []fs.DirEntry
	var fixedArray = []string{"base_name-1", "base_name-1.1", "base_name-2", "base_name-2.1"}
    for _, name := range fixedArray {
        entries = append(entries, fakeDirEntry{name: name})
    }
    return entries, nil
}


func TestGetLatestApp(t *testing.T) {
	readDirFunc = fakeReadDir
	appFolder := "test_folder"
	result, _ := GetLatestApp(appFolder, "base_name")

    // 使用 assert.Equal 来检查相等性
    assert.Equal(t, filepath.Join(filepath.Join(GetAppRoot(), appFolder), "base_name-2.1"), result)
}
