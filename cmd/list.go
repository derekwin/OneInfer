package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var modelListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all added models",
	Run: func(cmd *cobra.Command, args []string) {
		err := listModels()
		if err != nil {
			fmt.Println("Error:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(modelListCmd)
}

// listModels 列出所有保存的模型
func listModels() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	modelDir := filepath.Join(homeDir, ".oneinfer", "models")
	metaPath := filepath.Join(modelDir, "models.json")

	// 读取现有的 models.json
	if _, err := os.Stat(metaPath); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No models found.")
			return nil
		}
		return err
	}

	// 解析 models.json 文件
	var models []map[string]string
	file, err := os.ReadFile(metaPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &models)
	if err != nil {
		return err
	}

	// 输出模型列表
	if len(models) == 0 {
		fmt.Println("No models found.")
		return nil
	}

	// 打印模型信息
	fmt.Printf("%-30s %-20s %-50s\n", "Model Name", "Platform", "Path")
	fmt.Println(strings.Repeat("-", 100))
	for _, model := range models {
		// 打印每个模型的名称、平台和路径
		fmt.Printf("%-30s %-20s %-50s\n", model["name"], model["platform"], model["path"])
	}

	return nil
}
