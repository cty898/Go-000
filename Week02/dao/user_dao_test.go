package dao

import (
	"testing"
)

func TestUserDAOImpl_Save(t *testing.T) {
	userDAO := &UserDAOImpl{}

	err := InitMysql("127.0.0.1", "3306", "root", "123456", "user")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	user := &UserEntity{
		Username: "cty",
		Password: "123456",
		Email:    "cty@mail.com",
	}

	err = userDAO.Save(user)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("new User ID is %d", user.ID)

}

func TestUserDAOImpl_SelectByEmail(t *testing.T) {

	userDAO := &UserDAOImpl{}

	err := InitMysql("127.0.0.1", "3306", "root", "123456", "user")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	user, err := userDAO.SelectByEmail("cty@mail.com")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("result username is %s", user.Username)

}
