package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// ps 命令
var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "List all running models",
	Run: func(cmd *cobra.Command, args []string) {
		listRunningModels()
	},
}

func init() {
	rootCmd.AddCommand(psCmd)
}

// 获取所有运行的模型
func listRunningModels() {
	resp, err := http.Get("http://127.0.0.1:9090/models")
	if err != nil {
		fmt.Println("Error: Unable to connect to oneinfer service.")
		os.Exit(1)
	}
	defer resp.Body.Close()

	// 处理 HTTP 错误码
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: Server returned status %d\n", resp.StatusCode)
		os.Exit(1)
	}

	// 读取并打印响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error: Unable to read the response body.")
		os.Exit(1)
	}

	// 解析 JSON 响应
	var models []ModelProcessStatus
	if err := json.Unmarshal(body, &models); err != nil {
		fmt.Println("Error: Failed to parse response:", err)
		os.Exit(1)
	}

	// 没有运行的模型
	if len(models) == 0 {
		fmt.Println("No running models found.")
		return
	}

	// 输出表头
	fmt.Printf("%-10s %-20s %-15s %-6s %-10s\n", "PID", "Model", "Host", "Port", "Status")
	fmt.Println("--------------------------------------------------------------")

	// 输出每个模型的信息
	for _, model := range models {
		fmt.Printf("%-10d %-20s %-15s %-6d %-10s\n", model.ID, model.Model, model.Host, model.Port, model.Status)
	}
}
