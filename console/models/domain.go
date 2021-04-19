package models

import "fmt"

type Domain struct {
	ID        int       `gorm:"column:id"`
	Tag       string    `gorm:"column:tag"`
	Domains   string    `gorm:"column:domains"`
	Enable    int       `gorm:"column:enable"`
}

const DOMAIN_TABLE_NAME = "domains"

func DomainFind(tag string) *Domain {
	var out Domain
	db := orm.Table(DOMAIN_TABLE_NAME)
	db.Where("tag = ?", tag).First(&out)
	if out.ID == 0 {
		return nil
	}
	return &out
}

func DomainGet() []Domain {
	var out []Domain
	db := orm.Table(DOMAIN_TABLE_NAME)
	db.Find(&out)
	return out
}

func DoaminUpdate(tag string, update func(u *Domain) ) error {
	var info Domain
	db := orm.Table(DOMAIN_TABLE_NAME)
	db.Where("tag = ?", tag).First(&info)
	if info.ID == 0 {
		return fmt.Errorf("%s not exist", tag)
	}
	update(&info)
	db.Where("id = ?", info.ID).UpdateColumns(info)
	return nil
}

