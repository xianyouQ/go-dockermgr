package models

import (
	"errors"
	"log"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"github.com/xianyouQ/go-dockermgr/utils"
)

//用户表
type User struct {
	Id            int64
	Username      string         `orm:"unique;size(32)" form:"Username"  valid:"Required;MaxSize(20);MinSize(6)"`
	Password      string         `orm:"size(32)" form:"Password" valid:"Required;MaxSize(20);MinSize(6)"`
	Repassword    string         `orm:"-" form:"Repassword" valid:"Required"`
	OldPassword   string         `orm:"-" form:"Oldpassword"`
	Email         string         `orm:"size(32)" form:"Email" valid:"Email"`
	Lastlogintime time.Time      `orm:"null;type(datetime)" form:"-"`
	Createtime    time.Time      `orm:"type(datetime);auto_now_add" `
	ServiceAuths  []*ServiceAuth `orm:"rel(m2m);rel_through(github.com/xianyouQ/go-dockermgr/models.ServiceAuthUser)"`
}

func (u *User) TableName() string {
	return beego.AppConfig.String("rbac_user_table")
}

func (u *User) Valid(v *validation.Validation) {
	if u.Password != u.Repassword {
		v.SetError("Repassword", "两次输入的密码不一样")
	}
}

//验证用户信息
func checkUser(u *User) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(&u)
	if !b {
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
			return errors.New(err.Message)
		}
	}
	return nil
}

func init() {
	orm.RegisterModel(new(User))
}

/************************************************************/

//get user list
func Getuserlist(username string, page int64, page_size int64, sort string) ([]*User, error) {
	o := orm.NewOrm()
	var users []*User
	var err error
	qs := o.QueryTable(beego.AppConfig.String("rbac_user_table"))
	var offset int64
	if page <= 1 {
		offset = 0
	} else {
		offset = (page - 1) * page_size
	}
	if username != "" {
		_, err = qs.Filter("Username__istartswith", username).Limit(page_size, offset).OrderBy(sort).All(&users)
	} else {
		_, err = qs.Limit(page_size, offset).OrderBy(sort).All(&users)
	}
	return users, err
}

//添加用户
func AddUser(o orm.Ormer, u *User) (int64, error) {
	var err error
	var id int64
	if err = checkUser(u); err != nil {
		return 0, err
	}
	u.Password = utils.Strtomd5(u.Password)
	id, err = o.Insert(u)
	if err != nil && u.Username == beego.AppConfig.String("rbac_admin_user") {
		return id, err
	}
	var baseRole *Role
	baseRole, err = QueryRole("BASE")
	if err != nil {
		return id, err
	}
	users := make([]*User, 0, 1)
	users = append(users, u)
	err = AddUserAuth(o, users, baseRole, nil)
	return id, err
}

//更新用户
func UpdateUser(o orm.Ormer, u *User, updatecols ...string) (int64, error) {
	if err := checkUser(u); err != nil {
		return 0, err
	}
	num, err := o.Update(u, updatecols...)
	return num, err
}

func DelUserById(o orm.Ormer, Id int64) (int64, error) {
	status, err := o.Delete(&User{Id: Id})
	return status, err
}

func GetUserByUsername(username string) (user User) {
	user = User{Username: username}
	o := orm.NewOrm()
	o.Read(&user, "Username")
	return user
}
