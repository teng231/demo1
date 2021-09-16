package db

import (
	"errors"
	"fmt"
	"log"
	"time"

	"git.urbox.vn/urbackend/user/utils"
	_ "github.com/lib/pq"
	"github.com/teng231/demo1/pb"
	"xorm.io/xorm"
)

type DB struct {
	engine *xorm.Engine
}

// ConnectDb open connection to db
func (d *DB) ConnectDb(sqlPath, dbName, sslmode string) error {
	sqlConnStr := fmt.Sprintf("%s/%s?sslmode=%s", sqlPath, dbName, sslmode)
	log.Print(sqlConnStr)
	engine, err := xorm.NewEngine("postgres", sqlConnStr)
	if err != nil {
		return err
	}
	tick := time.NewTicker(15 * time.Minute)
	go func(engine *xorm.Engine) {
		for {
			select {
			case <-tick.C:
				if err := engine.Ping(); err != nil {
					log.Print("sql can not ping")
				}
			}
		}
	}(engine)
	log.Print("Connected to: ", sqlConnStr)
	d.engine = engine
	d.engine.ShowSQL(false)
	return err
}

func (d *DB) listUsersQuery(rq *pb.UserRequest) *xorm.Session {
	ss := d.engine.Table("user")
	if rq.GetUsername() != "" {
		ss.And("username = ?", rq.GetUsername())
	}
	if rq.GetFullname() != "" {
		ss.And("fullname like ?", "%"+rq.GetFullname()+"%")
	}
	if rq.GetPhone() != "" {
		ss.And("phone = ?", rq.GetPhone())
	}
	if rq.GetId() != 0 {
		ss.And("id = ?", rq.GetId())
	}
	if len(rq.GetIds()) != 0 {
		ss.In("id", rq.GetIds())
	}
	if len(rq.GetNotIds()) != 0 {
		ss.NotIn("id", rq.GetNotIds())
	}
	if rq.GetState() != 0 {
		ss.And("state = ?", rq.GetState())
	}
	return ss
}

// ListUsers ...
func (d *DB) ListUsers(rq *pb.UserRequest) ([]*pb.User, error) {

	ss := d.listUsersQuery(rq)
	if rq.GetLimit() != 0 {
		ss.Limit(int(rq.GetLimit()), int(rq.GetOffset()*rq.GetLimit()))
	}
	users := make([]*pb.User, 0)
	err := ss.Desc("id").Find(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// CountUsers ...
func (d *DB) CountUsers(rq *pb.UserRequest) (int64, error) {
	ss := d.listUsersQuery(rq)
	return ss.Count()
}

// FindUser get single user - oke
func (d *DB) FindUser(rq *pb.UserRequest) (*pb.User, error) {
	user := &pb.User{
		Id:       rq.GetId(),
		Username: rq.GetUsername(),
		Email:    rq.GetEmail(),
		Phone:    rq.GetPhone(),
	}
	ishas, err := d.engine.Desc("id").Get(user)
	if err != nil {
		return nil, err
	}
	if !ishas {
		return nil, errors.New(utils.E_error_not_found)
	}
	return user, nil
}
