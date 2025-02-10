package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var modelAddCmd = &cobra.Command{
	Use:   "add <model_name> <platform_name_or_local> [file_pattern]",
	Short: "Add a model by specifying either a local path or platform and model name, with an optional file pattern to limit the download",
	Args:  cobra.MinimumNArgs(2), // 至少两个参数，平台和模型名，第3个是可选的文件模式
	Run: func(cmd *cobra.Command, args []string) {
		modelName := args[0]
		platformOrPath := args[1]
		var filePattern string

		// 可选：检查是否有 file_pattern 参数
		if len(args) > 2 {
			filePattern = args[2]
		}

		err := addModel(modelName, platformOrPath, filePattern)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Printf("Model '%s' added successfully!\n", modelName)
	},
}

func init() {
	rootCmd.AddCommand(modelAddCmd)
}

// addModel 根据平台名或本地路径添加模型
func addModel(name, platformOrPath, filePattern string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	modelDir := filepath.Join(homeDir, ".oneinfer", "models")
	if err := os.MkdirAll(modelDir, 0755); err != nil {
		return err
	}

	var destPath string
	var errCopy error

	// 如果平台是 "local"，则提示用户输入本地文件路径
	if platformOrPath == "local" {
		var localPath string
		fmt.Print("Enter the local file path for the model: ")
		_, err := fmt.Scanln(&localPath)
		if err != nil {
			return fmt.Errorf("invalid input: %v", err)
		}

		// 检查本地路径是否有效
		if fileInfo, err := os.Stat(localPath); err == nil && !fileInfo.IsDir() {
			// 创建一个以模型名称为名的文件夹
			modelFolder := filepath.Join(modelDir, name)
			if err := os.MkdirAll(modelFolder, 0755); err != nil {
				return fmt.Errorf("failed to create model folder: %v", err)
			}

			// 设置目标路径为该文件夹中的模型文件
			destfilePath := filepath.Join(modelFolder, name+filepath.Ext(localPath))
			errCopy = copyFile(localPath, destfilePath)

			// 后续模型保存元数据
			name = name + filepath.Ext(localPath)
			destPath = modelFolder
		} else {
			return fmt.Errorf("invalid local model file path")
		}
	} else {
		// 如果是远程平台，使用 os/exec 调用 Python 下载模型
		destPath = modelDir
		errCopy = downloadModelWithPython(platformOrPath, name, modelDir, filePattern)
		name = filepath.Join(name, filePattern)
	}

	if errCopy != nil {
		return errCopy
	}

	// 保存模型的元数据
	metaPath := filepath.Join(modelDir, "models.json")
	err = saveModelMetadata(metaPath, name, platformOrPath, destPath)
	if err != nil {
		return err
	}

	return nil
}

// downloadModelWithPython 使用 os/exec 调用 Python 下载模型
func downloadModelWithPython(platform, modelName, destPath, filePattern string) error {
	// 创建 Python 脚本
	pythonScript := fmt.Sprintf(`
import subprocess
import sys
import importlib

# 安装 pip（如果没有安装）
def install_pip():
    try:
        subprocess.check_call([sys.executable, "-m", "ensurepip", "--upgrade"])
    except subprocess.CalledProcessError:
        print("Failed to install pip. Please install pip manually.")
        # sys.exit(1) # just skip

# 安装所需的库
def install_libraries(libraries):
    for lib in libraries:
        try:
            importlib.import_module(lib)
        except ImportError:
            print(f"Installing {lib}...")
            subprocess.check_call([sys.executable, "-m", "pip", "install", lib])

# 检查并安装库
def check_and_install_libraries():
    install_pip()
    install_libraries(['huggingface_hub', 'modelscope'])

# 下载模型
def download_model(platform, model_name, dest_path, file_pattern):
    if platform == 'huggingface':
        from huggingface_hub import snapshot_download
        snapshot_download(repo_id=model_name, cache_dir=dest_path, allow_patterns=file_pattern)
    elif platform == 'modelscope':
        from modelscope import snapshot_download
        snapshot_download(model_name, cache_dir=dest_path, allow_file_pattern=file_pattern)
    else:
        print(f"Unsupported platform: {platform}")
        return None

check_and_install_libraries()
download_model('%s', '%s', '%s', '%s')
`, platform, modelName, destPath, filePattern)

	// 使用 os/exec 执行 Python 脚本
	cmd := exec.Command("python3", "-c", pythonScript)

	// 设置输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error creating stdout pipe: %v", err)
	}

	// 设置错误输出管道
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("error creating stderr pipe: %v", err)
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting command: %v", err)
	}

	// 实时读取标准输出
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text()) // 打印 Python 脚本的输出，包含进度信息
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading stdout:", err)
		}
	}()

	// 实时读取标准错误输出
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Println("Error:", scanner.Text()) // 打印 Python 脚本的错误输出
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading stderr:", err)
		}
	}()

	// 等待命令完成
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("error downloading model using Python: %v", err)
	}

	return nil
}

// copyFile 复制本地文件
func copyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}

// saveModelMetadata 保存模型的元数据到 models.json
func saveModelMetadata(metaPath, name, platform, path string) error {
	var models []map[string]string

	// 读取现有的 models.json
	if _, err := os.Stat(metaPath); err == nil {
		file, err := os.ReadFile(metaPath)
		if err == nil {
			json.Unmarshal(file, &models)
		}
	}

	// 检查模型是否已经存在
	for _, model := range models {
		if model["name"] == name && model["platform"] == platform {
			// 如果模型已存在，则跳过添加
			fmt.Printf("Model '%s' already exists in the metadata.\n", name)
			return nil
		}
	}

	// 添加新模型
	path = filepath.Join(path, name)
	models = append(models, map[string]string{"name": name, "platform": platform, "path": path})
	data, err := json.MarshalIndent(models, "", "  ")
	if err != nil {
		return err
	}

	// 保存到文件
	return os.WriteFile(metaPath, data, 0644)
}
