package pb

import (
	"fmt"
	"os/exec"

	"golang.org/x/xerrors"
)

// CompileGRPC compiles protobuf file to gRPC service protobuf definition
func CompileGRPC(protocPath, protoFilePath string, protoFileIncludePath []string) error {
	// build args which includes proto file include path
	args := []string{}
	for _, path := range protoFileIncludePath {
		arg := fmt.Sprintf("--proto_path=%v", path)
		args = append(args, arg)
	}
	args = append(args, "--go_out=.", "--go_opt=paths=source_relative", "--go-grpc_out=.", "--go-grpc_opt=paths=source_relative")
	args = append(args, protoFilePath)

	// execute protoc to generate .pb.go file
	shellCmd := exec.Command(protocPath, args...)
	if err := shellCmd.Run(); err != nil {
		return xerrors.Errorf("%v %v failed! error:%v", protocPath, args, err.Error())
	}
	return nil
}
