make clean
make USE_CUDA=1 # only llama.cpp with cuda
# make USE_CUDA=1 ONNX_SERVER=yes # with onnx backend
sudo bash uninstall.sh
sudo bash install.sh