package UserService

import (
	"context"

	pb "github.com/DarkMetrix/gofra/test/proto"
)

func (h UserServiceImpl) AddUser (ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	resp := new(pb.AddUserResponse)
	resp.Session = new(pb.Session)
	resp.Result = new(pb.Result)

	//Set session
	resp.Session = req.Session

	//Set result
	resp.GetResult().Code = 0
	resp.GetResult().Message = "Success"

	return resp, nil
}
