# protoc-generate-go-demo

##  第一版
cd template/protoc-gen-go-demo

## 编译和安装和执行
 go build .
##
 mv protoc-gen-go-demo.exe ../../../../bin/   ----- 移动到GO_BIN目录中
## 
cd ../../
 protoc -I ./protos --go_out ./protos --go_opt=paths=source_relative --go-demo_out ./protos --go-demo_opt=paths=source_relative protos/v1.proto
 