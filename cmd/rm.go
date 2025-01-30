package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var modelRemoveCmd = &cobra.Command{
	Use:     "remove <model_name>",
	Aliases: []string{"rm"},
	Short:   "Remove a model by name",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		modelName := args[0]
		err := removeModel(modelName)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Printf("Model '%s' removed successfully!\n", modelName)
	},
}

func init() {
	rootCmd.AddCommand(modelRemoveCmd)
}

// removeModel 根据模型名删除模型
func removeModel(name string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	modelDir := filepath.Join(homeDir, ".oneinfer", "models")
	metaPath := filepath.Join(modelDir, "models.json")

	// 1. 从 models.json 文件中删除对应的模型元数据
	var models []map[string]string
	if _, err := os.Stat(metaPath); err == nil {
		file, err := os.ReadFile(metaPath)
		if err != nil {
			return err
		}

		// 解析现有的 models.json
		err = json.Unmarshal(file, &models)
		if err != nil {
			return err
		}

		// 找到并删除指定模型
		var updatedModels []map[string]string
		for _, model := range models {
			if model["name"] != name {
				updatedModels = append(updatedModels, model)
			}
		}

		// 如果模型未找到
		if len(updatedModels) == len(models) {
			return fmt.Errorf("model '%s' not found", name)
		}

		// 保存更新后的 models.json
		data, err := json.MarshalIndent(updatedModels, "", "  ")
		if err != nil {
			return err
		}

		err = os.WriteFile(metaPath, data, 0644)
		if err != nil {
			return err
		}
	}

	// 2. 删除模型文件
	modelFilePath := filepath.Join(modelDir, name)
	if err := os.RemoveAll(modelFilePath); err != nil {
		return fmt.Errorf("failed to delete model file: %v", err)
	}

	return nil
}
