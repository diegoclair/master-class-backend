package gapi

import (
	db "github.com/diegoclair/master-class-backend/db/sqlc"
	"github.com/diegoclair/master-class-backend/proto/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(u db.User) *pb.User {
	return &pb.User{
		Username:          u.Username,
		FullName:          u.FullName,
		Email:             u.Email,
		PasswordChangedAt: timestamppb.New(u.PasswordChangedAt),
		CreateAt:          timestamppb.New(u.CreatedAt),
	}
}
