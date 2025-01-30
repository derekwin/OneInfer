# OneInfer
一站式推理模型管理工具

OneInfer 是一个命令行工具，用于管理和提供多种类型的机器学习模型，支持本地和远程模型的无缝集成。它支持从 Hugging Face、ModelScope 等平台下载模型，并使用 `llama.cpp` 后端提供模型服务。

## 与 Ollama 的区别

与 Ollama 相比，OneInfer 提供了更多的灵活性和选择：

1. **更广泛的模型平台支持**：OneInfer 支持从多个平台下载模型，包括 Hugging Face 和 ModelScope，而 Ollama 限于其自有平台。
2. **支持多种推理后端**：OneInfer 将支持多种推理后端，不仅仅是语言模型，还包括视觉模型和其他非 LLM 模型，为用户提供更多选择和自由度。

通过 OneInfer，用户可以享受便捷的本地部署，同时能够享受更多的平台和模型选择，提供更丰富和灵活的体验。

## 开发路线
- [x] 类似 Ollama 的模型管理
- [x] 支持使用 Hugging Face 或 ModelScope 中的预训练模型
- [x] 支持通过 `llama.cpp` 后端提供 LLM 模型服务（.gguf 格式）
- [ ] 支持更多推理后端
- [ ] 支持更多类型的模型
- [ ] 提供开箱即用的打包应用，用户无需编译，直接下载即可使用

## 系统要求
- Python3：用于从 Hugging Face 和 ModelScope 下载模型。
- Go 1.18+：用于构建 OneInfer。
- git：用于下载其他代码库。

## 构建与安装
要构建并安装 OneInfer：

1. 克隆仓库：
   ```bash
   git clone https://github.com/derekwin/oneinfer.git
   cd oneinfer
   ```
2. 运行 `make` 命令：
   ```bash
   make USE_CUDA=1
   ```

   针对不同的 GPU/后端配置：
   - **USE_BLAS=1**：适用于普通 CPU。
   - **USE_CUDA=1**：适用于 NVIDIA GPU。
   - **USE_MUSA=1**：适用于 Meta 的 AI 加速器。
   - **USE_HIP=1**：适用于 AMD GPU。
   - **USE_CANN=1**：适用于华为 Ascend AI 加速器。
   - **USE_Vulkan=1**：适用于支持 Vulkan 的 GPU，提供高效的并行计算。
   - **USE_Metal=1**：适用于 Apple 设备。
   - **USE_SYCL=1**：适用于多种异构设备（包括 CPU、GPU、FPGA），使用 oneAPI。

3. 安装/卸载二进制文件：
   ```bash
   sudo bash install.sh/uninstall.sh
   ```

或者可以直接运行 `bash allinnoe.sh`。

## 使用方法

### 添加模型
将模型添加到 OneInfer。可以从 ModelScope、Hugging Face 或本地文件添加。
（您可以从 Hugging Face 或 ModelScope 网站获取仓库名称和文件名。）

```bash
oneinfer add <model_repo> <platform_name> <file_name>
```

#### 从 ModelScope 下载模型
例如，从 ModelScope 下载 `DeepSeek` 模型：

```bash
# deepseek r1 from unsloth
oneinfer add unsloth/DeepSeek-R1-Distill-Qwen-32B-GGUF modelscope DeepSeek-R1-Distill-Qwen-32B-Q5_K_M.gguf
oneinfer add unsloth/DeepSeek-R1-Distill-Qwen-7B-GGUF modelscope DeepSeek-R1-Distill-Qwen-7B-Q4_K_M.gguf
```

#### 从 Hugging Face 下载模型
例如，从 Hugging Face 下载模型：

```bash
oneinfer add RepoId huggingface modelname
```

#### 添加本地模型
例如，添加本地模型文件：

```bash
oneinfer add localmodelname local 
# 然后输入文件路径
./test/fakemodel.bin
```

### 列出模型
列出所有已添加到 OneInfer 的可用模型。

```bash
oneinfer ls
```

### 删除模型
通过模型名称删除特定模型。

```bash
oneinfer rm <model_name>
```

## 作为服务器运行
首先运行 OneInfer 作为后台服务器以管理模型服务：

```bash
nohup oneinfer serve &
```

这将启动 OneInfer 服务器，在后台管理模型服务。

## 作为客户端管理

### 启动模型
通过指定模型名称启动特定模型。您还可以定义模型服务器的主机和端口。

```bash
oneinfer run modelname [-p (默认值 8080)] [-h (默认值 127.0.0.1)]
```

例如：

```bash
oneinfer run -m DeepSeek-R1-Distill-Qwen-7B-Q4_K_M.gguf
```

这将调用 OneInfer 服务器并启动模型服务。

### 查看所有运行中的模型状态
查看所有正在运行的模型的状态：

```bash
oneinfer ps
```

这将列出当前运行的模型及其状态。

### 停止模型
通过模型的唯一标识符（UID）停止运行中的模型：

```bash
oneinfer stop <model_uid>
```

### 停止服务器
停止整个 OneInfer 服务器：

```bash
oneinfer stop serve
```

这将停止服务器及所有运行中的模型。

---

## 故障排除

- 如果遇到模型下载问题，请确保已正确安装并配置 Python 3，用于 ModelScope 和 Hugging Face 集成。
- 如果无法启动模型，请检查端口是否被占用或是否缺少任何依赖。

有关每个命令的详细帮助，可以使用 `--help` 标志：
```bash
oneinfer --help
```

---

## 许可
OneInfer 是一个开源软件，采用 MIT 许可证。