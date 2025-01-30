INSTALL_DIR="/usr/local/bin"
LLAMA_SERVER_DIR="/usr/local/oneinfer"

# 复制oneinfer到安装目录
rm $INSTALL_DIR/oneinfer

# 复制llama server到对应目录
rm -r $LLAMA_SERVER_DIR/*

echo "卸载成功"