# 在摩尔线程GPU服务器上使用OneInfer下载并部署DeepSeek模型

随着大语言模型（LLM）的广泛应用，如何高效地在本地部署这些模型成为了一个重要的问题。OneInfer为用户提供了便捷的本地模型部署和管理方案。本文将介绍如何在摩尔线程的GPU服务器上使用OneInfer下载并部署DeepSeek模型，帮助你在本地环境中快速启动和运行大语言模型。

## 一、环境准备

在开始之前，确保你的摩尔线程GPU服务器已经安装了以下必要的软件和驱动程序：

1. **安装摩尔线程驱动**

   摩尔线程的GPU需要安装相应的驱动程序才能正常工作。你需要从摩尔线程的官网下载最新的驱动包，并按照官方文档的步骤依次安装MT-Linux-driver和MUSA_Toolkits。确保安装过程中没有错误提示。

2. **安装依赖**

   在安装驱动后，还需要安装一些必要的依赖项。运行以下命令来安装这些依赖：
   
   ```bash
   sudo apt-get update
   sudo apt-get install -y build-essential cmake git g++ python3 python3-pip golang
   ```
   
   这些依赖项将为后续的编译和运行提供支持。

## 二、安装OneInfer

OneInfer是一个轻量级的CLI工具，支持从多个平台（如Hugging Face和ModelScope）下载模型，并通过llama.cpp后端进行优化的CPU/GPU模型推理。以下是安装OneInfer的步骤：

1. **克隆OneInfer仓库**
   
   首先，克隆OneInfer的GitHub仓库：
   
   ```bash
   git clone https://github.com/derekwin/oneinfer.git
   cd oneinfer
   ```

2. **编译OneInfer**
   
   在编译OneInfer时，需要指定使用摩尔线程的GPU后端。运行以下命令进行编译：
   
   ```bash
   make USE_MUSA=1
   ```
   
   这将启用对摩尔线程GPU的支持。

3. **安装OneInfer**
   
   编译完成后，运行以下命令安装OneInfer：
   
   ```bash
   sudo bash install.sh
   ```
   
   安装完成后，你可以通过`oneinfer`命令来管理模型的部署和运行。

## 三、使用OneInfer接口下载DeepSeek模型

OneInfer提供了便捷的接口来下载和管理模型，支持从ModelScope和Hugging Face等平台直接下载模型。以下是使用OneInfer接口下载DeepSeek模型的步骤：

1. **下载DeepSeek模型**
   
   使用OneInfer的add命令从ModelScope下载DeepSeek模型。运行以下命令：
   
   ```bash
   oneinfer add unsloth/DeepSeek-R1-Distill-Qwen-32B-GGUF modelscope DeepSeek-R1-Distill-Qwen-32B-Q5_K_M.gguf
   ```
   
   这条命令会自动从ModelScope平台下载DeepSeek-R1-Distill-Qwen-32B-Q5_K_M.gguf模型文件，并将其添加到OneInfer的模型管理中。

2. **查看已下载的模型**
   
   运行以下命令查看已添加到OneInfer的模型列表：
   
   ```bash
   oneinfer ls
   ```
   
   你将看到刚刚下载的DeepSeek模型出现在列表中。

## 四、部署并运行DeepSeek模型

1. **启动OneInfer服务器**
   
   在部署模型之前，需要启动OneInfer服务器来管理模型服务：
   
   ```bash
   nohup oneinfer serve &
   ```
   
   这将启动OneInfer服务器，使其在后台运行。

2. **启动DeepSeek模型**
   
   使用以下命令启动DeepSeek模型：
   
   ```bash
   oneinfer run DeepSeek-R1-Distill-Qwen-32B-Q5_K_M.gguf -p 8080 -H 0.0.0.0
   ```
   
   这将启动模型服务，并将其绑定到指定的端口（默认为8080）和主机地址（默认为0.0.0.0，表示允许所有IP访问）。

3. **查看运行状态**
   
   运行以下命令查看所有正在运行的模型及其状态：
   
   ```bash
   oneinfer ps
   ```
   
   你将看到DeepSeek模型的运行状态。

4. **访问模型**
   
   启动服务后，你可以通过浏览器访问`http://<your-server-ip>:8080`来使用DeepSeek模型进行交互。你也可以通过API调用来使用模型。

## 五、停止服务

如果需要停止运行的模型或关闭OneInfer服务器，可以使用以下命令：

1. **停止运行的模型**
   
   通过模型的唯一标识符（UID）停止运行的模型：
   
   ```bash
   oneinfer stop <model_uid>
   ```

2. **停止OneInfer服务器**
   
   停止整个OneInfer服务器及其所有运行的模型：
   
   ```bash
   oneinfer stop serve
   ```

## 六、总结

通过上述步骤，你可以在摩尔线程的GPU服务器上使用OneInfer接口轻松下载并部署DeepSeek模型。OneInfer提供了便捷的模型管理功能，支持多种平台和后端，使得本地部署和运行大语言模型变得更加简单。希望本文能帮助你在本地环境中快速启动和运行DeepSeek模型，享受大语言模型带来的强大功能。