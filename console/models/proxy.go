package models

import "fmt"

const (
	MODE_LOCAL  = "local"
	MODE_REMOTE = "remote"
	MODE_DOMAIN = "domain"
)

type Proxy struct {
	ID        int       `gorm:"column:id"`
	Tag       string    `gorm:"column:tag"`
	Enable    int       `gorm:"column:enable"`
	Iface     string    `gorm:"column:interface"`
	Port      int       `gorm:"column:port"`
	Auth      int       `gorm:"column:auth"`
	Protocal  string    `gorm:"column:protocal"`
	Remote    string    `gorm:"column:remote"`
	Mode      string    `gorm:"column:mode"`
	Status    string    `gorm:"column:status"`
}

const PROXY_TABLE_NAME = "proxys"

func ProxyFind(tag string) *Proxy {
	var out Proxy
	db := orm.Table(PROXY_TABLE_NAME)
	db.Where("tag = ?", tag).First(&out)
	if out.ID == 0 {
		return nil
	}
	return &out
}

func ProxyFindByID(id string) *Proxy {
	var out Proxy
	db := orm.Table(PROXY_TABLE_NAME)
	db.Where("id = ?", id).First(&out)
	if out.ID == 0 {
		return nil
	}
	return &out
}

func ProxyGet() []Proxy {
	var out []Proxy
	db := orm.Table(PROXY_TABLE_NAME)
	db.Find(&out)
	return out
}

func ProxyUpdate(tag string, update func(u *Proxy) ) error {
	var info Proxy
	db := orm.Table(PROXY_TABLE_NAME)
	db.Where("tag = ?", tag).First(&info)
	if info.ID == 0 {
		return fmt.Errorf("%s not exist", tag)
	}
	update(&info)
	db.Where("id = ?", info.ID).UpdateColumns(info)
	return nil
}


