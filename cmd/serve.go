package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

// 进程信息结构体
type ModelProcess struct {
	ID      int    `json:"id"`
	Model   string `json:"model"`
	Status  string `json:"status"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Command *exec.Cmd
}

type ModelProcessStatus struct {
	ID     int    `json:"id"`
	Model  string `json:"model"`
	Status string `json:"status"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
}

var (
	models   = make(map[int]*ModelProcess)
	modelMux = sync.Mutex{}
)

// 启动 REST API 服务
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the oneinfer model management service",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting oneinfer service on http://0.0.0.0:9090...")
		router := mux.NewRouter()
		router.HandleFunc("/models", listModelsHandler).Methods("GET")
		router.HandleFunc("/models", startModelHandler).Methods("POST")
		router.HandleFunc("/models/{id}", stopModelHandler).Methods("DELETE")
		router.HandleFunc("/stop", stopServerHandler).Methods("POST")
		router.HandleFunc("/health", healthCheckHandler).Methods("GET")

		log.Fatal(http.ListenAndServe(":9090", router))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

// 列出所有运行的模型
func listModelsHandler(w http.ResponseWriter, r *http.Request) {
	modelMux.Lock()
	defer modelMux.Unlock()

	modelList := make([]ModelProcess, 0, len(models))
	for _, mp := range models {
		modelList = append(modelList, *mp)
	}

	// 输出日志确认是否获取到模型数据
	fmt.Printf("Returning models: %v\n", modelList)

	// 创建一个新的切片，用来存储可序列化的模型数据
	serializableModels := make([]ModelProcessStatus, 0, len(models))
	for _, mp := range models {
		// 只取出模型的可序列化字段
		serializableModels = append(serializableModels, ModelProcessStatus{
			ID:     mp.ID,
			Model:  mp.Model,
			Host:   mp.Host,
			Port:   mp.Port,
			Status: mp.Status,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serializableModels)
}

// 启动一个模型进程（分离进程）
func startModelHandler(w http.ResponseWriter, r *http.Request) {
	modelMux.Lock()
	defer modelMux.Unlock()

	// 解析 JSON 请求
	var req struct {
		Model string `json:"model"`
		Host  string `json:"host"`
		Port  int    `json:"port"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 检查端口是否已被占用
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", req.Host, req.Port))
	if err != nil {
		http.Error(w, fmt.Sprintf("Port %d is already in use", req.Port), http.StatusConflict)
		return
	}
	listener.Close() // 关闭监听器，因为只是做占用检测

	// 运行 Llama.cpp 进程（独立进程）
	serverPath := "/usr/local/oneinfer/llama/llama-server"
	cmd := exec.Command(serverPath, "--host", req.Host, "--port", strconv.Itoa(req.Port), "--model", req.Model, "-ngl", "9999")

	// 分离进程，不让 serve 进程被阻塞
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		http.Error(w, "Failed to start model", http.StatusInternalServerError)
		return
	}

	// 记录进程信息
	modelProcess := &ModelProcess{
		ID:      cmd.Process.Pid,
		Model:   req.Model,
		Host:    req.Host,
		Port:    req.Port,
		Command: cmd,
	}
	models[cmd.Process.Pid] = modelProcess

	// 创建一个新的切片，用来存储可序列化的模型数据
	serializableModels := make([]ModelProcessStatus, 0, len(models))
	for _, mp := range models {
		// 只取出模型的可序列化字段
		serializableModels = append(serializableModels, ModelProcessStatus{
			ID:     mp.ID,
			Model:  mp.Model,
			Host:   mp.Host,
			Port:   mp.Port,
			Status: mp.Status,
		})
	}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(serializableModels)
}

// 停止模型进程
func stopModelHandler(w http.ResponseWriter, r *http.Request) {
	modelMux.Lock()
	defer modelMux.Unlock()

	vars := mux.Vars(r)
	pid, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid process ID", http.StatusBadRequest)
		return
	}

	modelProcess, exists := models[pid]
	if !exists {
		http.Error(w, "Process not found", http.StatusNotFound)
		return
	}

	// 终止进程
	if err := syscall.Kill(-modelProcess.ID, syscall.SIGKILL); err != nil {
		http.Error(w, "Failed to stop process", http.StatusInternalServerError)
		return
	}

	delete(models, pid)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Stopped process %d\n", pid)
}

// 关闭 `serve` 并停止所有模型
func stopServerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Stopping oneinfer service...")

	// 停止所有模型进程
	modelMux.Lock()
	for pid, model := range models {
		fmt.Printf("Stopping model %s (PID %d)...\n", model.Model, pid)
		syscall.Kill(-pid, syscall.SIGKILL) // 终止进程组
		delete(models, pid)
	}
	modelMux.Unlock()

	fmt.Println("All models stopped. Exiting...")
	os.Exit(0)
}

// 健康检查
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
