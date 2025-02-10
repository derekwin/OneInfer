INSTALL_DIR="/usr/local/bin"
SERVER_DIR="/usr/local/oneinfer"

# 删除oneinfer到安装目录
rm $INSTALL_DIR/oneinfer

# 删除llama server到对应目录
rm -r $SERVER_DIR/*

echo "卸载成功"