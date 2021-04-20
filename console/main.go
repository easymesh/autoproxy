package main

import (
	"flag"
	"github.com/astaxie/beego/context"
	"github.com/GoAdminGroup/go-admin/modules/logger"
	"github.com/easymesh/autoproxy"
	"net/http"
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
)

var (
	Debug   bool
	Help    bool

	TlsCert    string
	TlsKey     string
	TlsEnable  bool

	Port    int
	Config  string
)

func init()  {
	flag.BoolVar(&Help, "help", false, "usage help")
	flag.BoolVar(&Debug, "debug", false, "debug")

	flag.IntVar(&Port, "port", 8000, "bind port for console")
	flag.StringVar(&Config, "config", "./config.json", "console config file")

	flag.BoolVar(&TlsEnable, "tls", false, "console with tls")
	flag.StringVar(&TlsCert, "cert", "", "certificate file")
	flag.StringVar(&TlsKey, "key", "", "private key file name")
}

func ProccessExit(eng *engine.Engine)  {
	eng.SqliteConnection().Close()
	pages.EnginFini()
	logger.Info("console shutdown")
	os.Exit(-1)
}

func main() {
	flag.Parse()
	if Help {
		flag.Usage()
		return
	}

	app := beego.NewApp()

	template.AddComp(chartjs.NewChart())

	beego.SetStaticPath("/uploads", "uploads")

	eng := engine.Default()
	eng.AddConfigFromJSON(Config)

	eng.AddGenerator("users", pages.UserTableGet)
	eng.AddGenerator("domains", pages.DomainTableGet)
	eng.AddGenerator("remotes", pages.RemoteTableGet)
	eng.AddGenerator("proxys", pages.ProxyTableGet)

	if err := eng.Use(app); err != nil {
		logger.Error(err.Error())
		panic(err)
	}
	//eng.HTML("GET", "/admin/", pages.GetDashBoard)
	app.Handlers.Any("/admin/", func(ctx *context.Context) {
		ctx.Redirect(http.StatusFound,"/admin/info/proxys")
	})

	app.Handlers.Any("/", func(ctx *context.Context) {
		ctx.Redirect(http.StatusFound,"/admin")
	})

	models.Init(eng.SqliteConnection())

	if TlsEnable {
		beego.BConfig.Listen.EnableHTTP  = false
		beego.BConfig.Listen.EnableHTTPS = true
		beego.BConfig.Listen.HTTPSAddr   = "0.0.0.0"
		beego.BConfig.Listen.HTTPSPort   = Port
		beego.BConfig.Listen.HTTPSCertFile = TlsCert
		beego.BConfig.Listen.HTTPSKeyFile  = TlsKey
	} else {
		beego.BConfig.Listen.HTTPAddr = "0.0.0.0"
		beego.BConfig.Listen.HTTPPort = Port
	}

	go func() {
		app.Run()
		logger.Info("proxy service shutdown")
		ProccessExit(eng)
	}()

	pages.DomainInit()
	pages.EnginInit()

	logger.Infof("proxy server %s start success", autoproxy.VersionGet())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	ProccessExit(eng)
}
