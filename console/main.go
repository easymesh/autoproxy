package main

import (
	"log"
	"os"
	"os/signal"

	_ "github.com/GoAdminGroup/go-admin/adapter/beego"             // web framework adapter
	_ "github.com/GoAdminGroup/go-admin/modules/db/drivers/sqlite" // sql driver
	_ "github.com/GoAdminGroup/themes/adminlte"                    // ui theme

	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/chartjs"
	"github.com/astaxie/beego"

	"github.com/easymesh/autoproxy/console/models"
	"github.com/easymesh/autoproxy/console/pages"
	"github.com/easymesh/autoproxy/console/tables"
)

func main() {
	startServer()
}

func startServer() {
	app := beego.NewApp()

	template.AddComp(chartjs.NewChart())

	eng := engine.Default()

	beego.SetStaticPath("/uploads", "uploads")

	if err := eng.AddConfigFromJSON("./config.json").
		AddGenerators(tables.Generators).
		Use(app); err != nil {
		panic(err)
	}

	eng.HTML("GET", "/admin", pages.GetDashBoard)
	eng.HTMLFile("GET", "/admin/hello", "./html/hello.tmpl", map[string]interface{}{
		"msg": "Hello world",
	})

	models.Init(eng.SqliteConnection())

	beego.BConfig.Listen.HTTPAddr = "127.0.0.1"
	beego.BConfig.Listen.HTTPPort = 80
	go app.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Print("closing database connection")
	eng.SqliteConnection().Close()
}
