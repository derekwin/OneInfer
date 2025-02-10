# `llama.cpp` 代码存放目录
LLAMA_DIR = $(abspath ./llama.cpp)
LLAMA_SERVER_BIN_PATH = $(abspath ./llama-server/)
# `onnxruntime-server` 代码存放目录
ONNX_DIR = $(abspath ./onnxruntime-server)
ONNX_SERVER_BIN_PATH = $(abspath ./onnx-server/)

# 目标架构
OS := $(shell uname -s)
ARCH := $(shell uname -m)

# Go 编译的最终二进制文件
ONEINFER_BIN = ./oneinfer

# LLAMACPP CMake选项
LLAMACPP_CMAKE_OPTS = -B$(LLAMA_DIR)/build -DCMAKE_BUILD_TYPE=Release
# LLAMACPP 默认后端选项
ifdef USE_BLAS
    LLAMACPP_CMAKE_OPTS += -DGGML_BLAS=ON
endif
ifdef USE_CUDA
    LLAMACPP_CMAKE_OPTS += -DGGML_CUDA=ON
endif
ifdef USE_MUSA
    LLAMACPP_CMAKE_OPTS += -DGGML_MUSA=ON
endif
ifdef USE_HIP
    LLAMACPP_CMAKE_OPTS += -DGGML_HIP=ON
endif
ifdef USE_CANN
    LLAMACPP_CMAKE_OPTS += -DGGML_CANN=ON
endif
ifdef USE_VULKAN
    LLAMACPP_CMAKE_OPTS += -DGGML_VULKAN=ON
endif
ifdef USE_METAL
    LLAMACPP_CMAKE_OPTS += -DGGML_METAL=ON
endif
ifdef USE_SYCL
    LLAMACPP_CMAKE_OPTS += -DGGML_SYCL=ON
endif

# ONNXRUNTIME_SERVER CMake选项

# 可选的编译 llama.cpp
LLAMA_CPP ?= yes
ONNX_SERVER ?= no

# 默认目标：编译 oneinfer
.PHONY: all clean build llama copy_so_libs clone_llama clone_onnx

all: build

# 1. 克隆 llama.cpp（如果不存在） 只有当 LLAMA_CPP 设置为 yes 时才执行
ifeq ($(LLAMA_CPP), yes)
clone_llama: 
	@if [ ! -d "$(LLAMA_DIR)" ]; then \
		echo "Cloning llama.cpp..."; \
		git clone https://github.com/ggerganov/llama.cpp $(LLAMA_DIR); \
	else \
		echo "llama.cpp already exists."; \
	fi
else
clone_llama:
	@echo "Skipping llama.cpp clone as LLAMA_CPP is set to no."
endif

# 2. 克隆 onnxruntime-server（如果不存在） 只有当 ONNX_SERVER 设置为 yes 时才执行
ifeq ($(ONNX_SERVER), yes)
clone_onnx:
	@if [ ! -d "$(ONNX_DIR)" ]; then \
		echo "Cloning onnxruntime-server..."; \
		git clone https://github.com/kibae/onnxruntime-server.git $(ONNX_DIR); \
	else \
		echo "onnxruntime-server already exists."; \
	fi
else
clone_onnx:
	@echo "Skipping onnxruntime-server clone as ONNX_SERVER is set to no."
endif

# 3. 编译 `llama.cpp server` 使用 CMake，只有当 LLAMA_CPP 设置为 yes 时才执行
ifeq ($(LLAMA_CPP), yes)
llama: clone_llama
	@if [ ! -f "$(LLAMA_SERVER_BIN_PATH)/llama-server" ]; then \
		echo "Compiling llama server with CMake..."; \
		cd $(LLAMA_DIR) && \
		mkdir -p build && \
		cd build && \
		cmake $(LLAMACPP_CMAKE_OPTS) .. && \
		make llama-server -j8; \
		mkdir -p $(LLAMA_SERVER_BIN_PATH); \
		cp $(LLAMA_DIR)/build/bin/llama-server $(LLAMA_SERVER_BIN_PATH); \
	else \
		echo "Llama server already compiled."; \
	fi
else
llama:
	@echo "Skipping llama server compilation as LLAMA_CPP is set to no."
endif

ifeq ($(ONNX_SERVER), yes)
prepare_onnx_server_env:
	@if [ "$(OS)" = "Linux" ]; then \
		if ldconfig -p | grep -q onnxruntime; then \
			echo "ONNX Runtime is already installed. Skipping setup."; \
		else \
			echo "Setting up ONNX Runtime on Linux..."; \
			cd $(ONNX_DIR) && bash download-onnxruntime-linux.sh; \
			if [ $$? -ne 0 ]; then \
				echo "Error: Failed to download ONNX Runtime. Please check your network connection." >&2; \
				exit 1; \
			fi; \
			echo "Adding /usr/local/onnxruntime/lib to /etc/ld.so.conf.d/onnxruntime.conf"; \
			echo "/usr/local/onnxruntime/lib" | sudo tee /etc/ld.so.conf.d/onnxruntime.conf > /dev/null; \
			sudo ldconfig; \
			echo "Installing dependencies..."; \
			sudo apt update && sudo apt install -y cmake pkg-config libboost-all-dev libssl-dev; \
		fi; \
	elif [ "$(OS)" = "Darwin" ]; then \
		if [ -d "/usr/local/include/onnxruntime" ] || [ -d "/opt/homebrew/include/onnxruntime" ]; then \
			echo "ONNX Runtime is already installed. Skipping setup."; \
		else \
			echo "Setting up ONNX Runtime on macOS..."; \
			brew install onnxruntime cmake boost openssl; \
		fi; \
	else \
		echo "Unsupported OS: $(OS)"; \
		exit 1; \
	fi
endif
# 编译 `onnxruntime-server` 使用 CMake，只有当 ONNX_SERVER 设置为 yes 时才执行
ifeq ($(ONNX_SERVER), yes)
onnx_server: clone_onnx prepare_onnx_server_env
	@if [ ! -f "$(ONNX_SERVER_BIN_PATH)/onnxruntime_server" ]; then \
		echo "Compiling onnxruntime server with CMake..."; \
		cd $(ONNX_DIR) && \
		cmake -B build -S . -DCMAKE_BUILD_TYPE=Release && \
		cmake --build build --parallel && \
		mkdir -p $(ONNX_SERVER_BIN_PATH) && \
		cp $(ONNX_DIR)/build/src/standalone/onnxruntime_server $(ONNX_SERVER_BIN_PATH); \
	else \
		echo "Onnxruntime server already compiled."; \
	fi
else
onnx_server:
	@echo "Skipping ONNXRUNTIME_SERVER compilation as ONNX_SERVER is set to no."
endif

# 4. 查找并复制所有生成的 `.so` 文件到 `$(SERVER_BIN)` 目录，只有当 LLAMA_CPP 设置为 yes 时才执行
ifeq ($(LLAMA_CPP), yes)
copy_so_libs: llama
	@echo "Copying .so libraries for llama..."
	@if [ -d "$(LLAMA_DIR)/build/bin" ]; then \
		mkdir -p $(LLAMA_SERVER_BIN_PATH); \
		cp $(LLAMA_DIR)/build/bin/*.so $(LLAMA_SERVER_BIN_PATH); \
	else \
		echo "No .so libraries found for llama."; \
	fi
else
copy_so_libs:
	@echo "Skipping copying .so libraries for llama as LLAMA_CPP is set to no."
endif

ifeq ($(ONNX_SERVER), yes)
copy_so_libs: onnx_server
	@echo "Copying .so libraries for onnxruntime-server..."
	@if [ -d "$(ONNX_DIR)/build/src" ]; then \
		mkdir -p $(ONNX_SERVER_BIN_PATH); \
		cp $(ONNX_DIR)/build/src/*.so $(ONNX_SERVER_BIN_PATH); \
	else \
		echo "No .so libraries found for onnxruntime-server."; \
	fi
else
copy_so_libs:
	@echo "Skipping copying .so libraries for onnxruntime-server as ONNX_SERVER is set to no."
endif


# 5. 编译 Go 项目
build: llama copy_so_libs onnx_server
	@if [ ! -f $(ONEINFER_BIN) ]; then \
		echo "Building oneinfer..."; \
		go build -o $(ONEINFER_BIN) ./main.go; \
	else \
		echo "oneinfer binary already exists."; \
	fi

# 6. 运行 `oneinfer`
run: build
	./$(ONEINFER_BIN)

# 7. 清理
clean:
	rm -rf $(LLAMA_DIR)/build $(LLAMA_SERVER_BIN_PATH) $(ONEINFER_BIN) llama-server
