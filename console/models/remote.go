package models

import "fmt"

type Remote struct {
	ID        int       `gorm:"column:id"`
	Tag       string    `gorm:"column:tag"`
	Enable    int       `gorm:"column:enable"`
	Address   string    `gorm:"column:address"`
	Auth      int       `gorm:"column:auth"`
	User      string    `gorm:"column:user"`
	Password  string    `gorm:"column:password"`
	Protocal  string    `gorm:"column:protocal"`
	Status    string    `gorm:"column:status"`
}

const REMOTE_TABLE_NAME = "remotes"

func RemoteFind(tag string) *Remote {
	var out Remote
	db := orm.Table(REMOTE_TABLE_NAME)
	db.Where("tag = ?", tag).First(&out)
	if out.ID == 0 {
		return nil
	}
	return &out
}

func RemoteGet() []Remote {
	var out []Remote
	db := orm.Table(REMOTE_TABLE_NAME)
	db.Find(&out)
	return out
}

func RemoteUpdate(tag string, update func(u *Remote) ) error {
	var info Remote
	db := orm.Table(REMOTE_TABLE_NAME)
	db.Where("tag = ?", tag).First(&info)
	if info.ID == 0 {
		return fmt.Errorf("%s not exist", tag)
	}
	update(&info)
	db.Where("id = ?", info.ID).UpdateColumns(info)
	return nil
}
