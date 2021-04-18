package models

import "fmt"

type User struct {
	ID        int       `gorm:"column:id"`
	User      string    `gorm:"column:user"`
	Password  string    `gorm:"column:password"`
	Enable    int       `gorm:"column:enable"`
	Flow      int       `gorm:"column:flow"`
	Online    int       `gorm:"column:online"`
}

const USER_TABLE_NAME = "users"

func UserFind(name string) *User {
	var out User
	db := orm.Table(USER_TABLE_NAME)
	db.Where("user = ?", name).First(&out)
	if out.ID == 0 {
		return nil
	}
	return &out
}

func UserGet() []User {
	var out []User
	db := orm.Table(USER_TABLE_NAME)
	db.Find(&out)
	return out
}

func UserUpdate(name string, update func(u *User) ) error {
	var info User
	db := orm.Table(USER_TABLE_NAME)
	db.Where("user = ?", name).First(&info)
	if info.ID == 0 {
		return fmt.Errorf("%s not exist", name)
	}
	update(&info)
	db.Where("id = ?", info.ID).UpdateColumns(info)
	return nil
}

