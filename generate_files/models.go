package generate_files

import (
	"os"
	"path/filepath"
)

//检查是否已存在models目录
func initModels(root string) error {
	_, err := os.Stat(filepath.Join(root, "models"))
	if err != nil {
		return os.MkdirAll(filepath.Join(root, "models"), 0777)
	}
	return nil
}
