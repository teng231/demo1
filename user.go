package main

import (
	"context"
	"errors"

	"github.com/teng231/demo1/pb"
)

func (d *Demo) ListUsers(ctx context.Context, in *pb.UserRequest) (*pb.Users, error) {
	// validate
	users, err := d.db.ListUsers(in)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return &pb.Users{}, nil
	}
	count, err := d.db.CountUsers(in)
	if err != nil {
		return nil, err
	}
	return &pb.Users{
		Users: users,
		Total: int32(count),
	}, nil
}

func (d *Demo) GetUser(ctx context.Context, in *pb.UserRequest) (*pb.User, error) {
	// validate
	user, err := d.db.FindUser(in)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (d *Demo) GetUserByUserName(ctx context.Context, in *pb.UserRequest) (*pb.User, error) {
	// validate
	if in.Username == "" {
		return nil, errors.New("not found username")
	}
	user, err := d.db.FindUser(in)
	if err != nil {
		return nil, err
	}
	return user, nil
}
