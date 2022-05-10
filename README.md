# protoc-generate-go-demo


## 编译和安装和执行
 go build .
##
 mv protoc-generate-go-demo.exe  protoc-gen-go-demo.exe
##
 mv protoc-gen-go-demo.exe ../../bin/   ----- 移动到GO_BIN目录中
## 
 protoc -I ./protos --go_out ./protos --go_opt=paths=source_relative --go-demo_out ./protos --go-demo_opt=paths=source_relative protos/v1.proto
 