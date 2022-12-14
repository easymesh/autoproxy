package main

import "github.com/easymesh/autoproxy/engin"

func init() {
	engin.SetVersion(VersionGet())
}

func VersionGet() string {
	return "v1.5.1"
}
