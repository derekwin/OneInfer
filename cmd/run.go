package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run <model_name>",
	Short: "Start a model through the oneinfer serve process",
	Args:  cobra.ExactArgs(1), // 确保 model_name 是必须的参数
	Run: func(cmd *cobra.Command, args []string) {
		// 获取命令行参数
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		modelName := args[0] // modelName 从 args 中获取

		// 默认值检查
		if host == "" {
			host = "127.0.0.1"
		}
		if port == 0 {
			port = 8080
		}

		// 获取模型路径
		modelPath, err := getModelPath(modelName)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// 构造请求数据
		requestBody, _ := json.Marshal(map[string]interface{}{
			"model": modelPath,
			"host":  host,
			"port":  port,
		})

		// 发送 REST 请求给 serve 进程
		resp, err := http.Post("http://127.0.0.1:9090/models", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			fmt.Println("Failed to request oneinfer serve:", err)
			return
		}
		defer resp.Body.Close()

		// 读取响应
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusCreated {
			fmt.Println("Error from server:", string(body))
			return
		}

		fmt.Println("Model started successfully:", string(body))
	},
}

func init() {
	// 添加命令行参数
	runCmd.Flags().StringP("host", "H", "", "IP address of the server (default is 127.0.0.1)")
	runCmd.Flags().IntP("port", "p", 8080, "Port number of the server (default is 8080)")

	// 添加 run 命令
	rootCmd.AddCommand(runCmd)
}

// getModelPath 获取指定模型的路径
func getModelPath(modelName string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	modelDir := filepath.Join(homeDir, ".oneinfer", "models")
	metaPath := filepath.Join(modelDir, "models.json")

	// 读取 models.json
	if _, err := os.Stat(metaPath); os.IsNotExist(err) {
		return "", fmt.Errorf("No models found.")
	}

	var models []map[string]string
	file, err := os.ReadFile(metaPath)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(file, &models)
	if err != nil {
		return "", err
	}

	// 如果没有模型
	if len(models) == 0 {
		return "", fmt.Errorf("No models found.")
	}

	// 查找指定名称的模型
	for _, model := range models {
		if model["name"] == modelName {
			return model["path"], nil
		}
	}

	return "", fmt.Errorf("Model '%s' not found.", modelName)
}
