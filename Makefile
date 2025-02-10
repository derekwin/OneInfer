# `llama.cpp` 代码存放目录
LLAMA_DIR = $(abspath ./llama.cpp)
SERVER_BIN_PATH = $(abspath ./llama-server/)

# 目标架构
OS := $(shell uname -s)
ARCH := $(shell uname -m)

# Go 编译的最终二进制文件
ONEINFER_BIN = ./oneinfer

# CMake选项
CMAKE_OPTS = -B$(LLAMA_DIR)/build -DCMAKE_BUILD_TYPE=Release

# 默认后端选项
ifdef USE_BLAS
    CMAKE_OPTS += -DGGML_BLAS=ON
endif

ifdef USE_CUDA
    CMAKE_OPTS += -DGGML_CUDA=ON
endif

ifdef USE_MUSA
    CMAKE_OPTS += -DGGML_MUSA=ON
endif

ifdef USE_HIP
    CMAKE_OPTS += -DGGML_HIP=ON
endif

ifdef USE_CANN
    CMAKE_OPTS += -DGGML_CANN=ON
endif

ifdef USE_VULKAN
    CMAKE_OPTS += -DGGML_VULKAN=ON
endif

ifdef USE_METAL
    CMAKE_OPTS += -DGGML_METAL=ON
endif

ifdef USE_SYCL
    CMAKE_OPTS += -DGGML_SYCL=ON
endif

# 默认目标：编译 oneinfer
.PHONY: all clean build llama copy_so_libs

all: build

# 1. 克隆 llama.cpp（如果不存在）
$(LLAMA_DIR):
	@if [ ! -d "$(LLAMA_DIR)" ]; then \
		echo "Cloning llama.cpp..."; \
		git clone https://github.com/ggerganov/llama.cpp $(LLAMA_DIR); \
	else \
		echo "llama.cpp already exists."; \
	fi

# 2. 编译 `llama.cpp server` 使用 CMake
llama: $(LLAMA_DIR)
	@if [ ! -f "$(SERVER_BIN)" ]; then \
		echo "Compiling llama server with CMake..."; \
		cd $(LLAMA_DIR) && \
		mkdir -p build && \
		cd build && \
		cmake $(CMAKE_OPTS) .. && \
		make llama-server -j8; \
		mkdir -p $(SERVER_BIN_PATH); \
		cp $(LLAMA_DIR)/build/bin/llama-server $(SERVER_BIN_PATH); \
	else \
		echo "Llama server already compiled."; \
	fi

# 3. 查找并复制所有生成的 `.so` 文件到 `$(SERVER_BIN)` 目录
copy_so_libs:
	@echo "Copying .so libraries..."
	@if [ -d "$(LLAMA_DIR)/build/bin" ]; then \
		mkdir -p $(SERVER_BIN_PATH); \
		cp $(LLAMA_DIR)/build/bin/*.so $(SERVER_BIN_PATH); \
	else \
		echo "No .so libraries found."; \
	fi

# 4. 编译 Go 项目
build: llama copy_so_libs
	@echo "Building oneinfer..."; \
	go build -o $(ONEINFER_BIN) ./main.go; \

# 5. 运行 `oneinfer`
run: build
	./$(ONEINFER_BIN)

# 6. 清理
clean:
	rm -rf $(LLAMA_DIR)/build $(SERVER_BIN) $(ONEINFER_BIN) llama-server
