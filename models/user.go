package models

import (
	"errors"
	// "strconv"
	// "time"

    "github.com/astaxie/beego/orm"
    "golang.org/x/crypto/bcrypt"
)

var (
	UserList map[string]*User
)

func init() {
}

type User struct {
    Id         int    `orm:"pk;unique;auto;column(user_id)"`
    Email      string
    // FamilyName string
    // GivenName  string
    Password  string    `json:"-"`
}

func (u *User) TableName() string {
    return "users"
}

func AddUser(u User) int {
	return 0
}

func GetUser(id int) (u *User, err error) {
	var user User
	o := orm.NewOrm()
	err = o.QueryTable("users").Filter("id", id).One(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetAllUsers() []*User {
	o := orm.NewOrm()
	var users []*User
	_, _ = o.QueryTable("users").All(&users)
	return users;
}

func UpdateUser(uid string, uu *User) (a *User, err error) {
	return nil, errors.New("User Not Exist")
}

func hash(s string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hash), err
}

func verify(hash, s string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(s))
}

func Login(email, password string) bool {
	for _, u := range UserList {
		if u.Email == email && u.Password == password {
			return true
		}
	}
	return false
}

func DeleteUser(uid string) {
	delete(UserList, uid)

}
