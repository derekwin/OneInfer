package cmd

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a model or the entire service",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Model ID or 'serve' is required")
		}

		// 如果参数是 "serve"，则停止整个服务
		if args[0] == "serve" {
			stopServer()
		} else {
			// 否则尝试停止指定模型
			modelID, err := strconv.Atoi(args[0])
			if err != nil {
				log.Fatalf("Invalid model ID: %v", err)
			}
			stopModel(modelID)
		}
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

// 停止指定模型的进程
func stopModel(modelID int) {
	url := fmt.Sprintf("http://localhost:9090/models/%d", modelID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Fatalf("Error creating DELETE request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending DELETE request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Printf("Model %d stopped successfully", modelID)
	} else {
		log.Printf("Failed to stop model %d, status code: %d", modelID, resp.StatusCode)
	}
}

// 停止整个服务
func stopServer() {
	url := "http://localhost:9090/stop"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		log.Fatalf("%v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Println("Server stopped successfully")
		os.Exit(0)
	} else {
		log.Printf("Failed to stop server, status code: %d", resp.StatusCode)
	}
}
