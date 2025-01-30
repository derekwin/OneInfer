#!/bin/bash

# 安装目录
INSTALL_DIR="/usr/local/bin"
LLAMA_SERVER_DIR="/usr/local/oneinfer/llama"

# 创建安装目录（如果不存在）
mkdir -p $INSTALL_DIR
mkdir -p $LLAMA_SERVER_DIR

# 复制oneinfer到安装目录
cp -f oneinfer $INSTALL_DIR

# 复制llama server到对应目录
cp -rf ./llama-server/* $LLAMA_SERVER_DIR

# 添加执行权限
chmod +x $INSTALL_DIR/oneinfer
chmod +x $LLAMA_SERVER_DIR/llama-server

# 更新系统环境变量，避免重复添加
if ! grep -q "$INSTALL_DIR" ~/.bashrc; then
    echo "export PATH=\$PATH:$INSTALL_DIR" >> ~/.bashrc
    echo "已更新PATH环境变量"
else
    echo "PATH环境变量已包含$INSTALL_DIR，跳过更新"
fi

# 刷新环境变量
source ~/.bashrc

echo "安装完成！oneinfer和llama server已安装至系统目录。"
