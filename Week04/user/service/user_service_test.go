package service

import (
	"context"
	"testing"

	"github.com/cty898/Go-000/Week04/user/dao"
	"github.com/cty898/Go-000/Week04/user/redis"
)

func TestUserServiceImpl_Login(t *testing.T) {

	err := dao.InitMysql("127.0.0.1", "3306", "root", "123456", "user")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = redis.InitRedis("127.0.0.1", "6379", "")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	userService := &UserServiceImpl{
		userDAO: &dao.UserDAOImpl{},
	}

	user, err := userService.Login(context.Background(), "cty@mail.com", "cty")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("user id is %d", user.ID)

}

func TestUserServiceImpl_Register(t *testing.T) {

	err := dao.InitMysql("127.0.0.1", "3306", "root", "123456", "user")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = redis.InitRedis("127.0.0.1", "6379", "")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	userService := &UserServiceImpl{
		userDAO: &dao.UserDAOImpl{},
	}

	user, err := userService.Register(context.Background(),
		&RegisterUserVO{
			Username: "cty",
			Password: "cty",
			Email:    "cty@mail.com",
		})

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("user id is %d", user.ID)

}
