package pages

import (
	"fmt"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/logger"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	editType "github.com/GoAdminGroup/go-admin/template/types/table"
	"github.com/easymesh/autoproxy/console/models"
	"github.com/easymesh/autoproxy/console/uitl"
	"strings"
	"sync"
	"time"
)

const (
	OPERATION_START   = "op_start"
	OPERATION_STOP    = "op_stop"
	OPERATION_RESTART = "op_restart"
)

func operationGet(path string) string {
	if strings.LastIndex(path, OPERATION_START) != -1 {
		return OPERATION_START
	}
	if strings.LastIndex(path, OPERATION_STOP) != -1 {
		return OPERATION_STOP
	}
	if strings.LastIndex(path, OPERATION_RESTART) != -1 {
		return OPERATION_RESTART
	}
	return ""
}

func enginProcess(ctx *context.Context) (success bool, msg string, data interface{}) {
	ids := strings.Split(ctx.FormValue("ids"), ",")
	var proxys []*models.Proxy
	for _, v := range ids {
		if add := models.ProxyFindByID(v) ; add != nil {
			proxys = append(proxys, add)
		}
	}
	if len(proxys) == 0 {
		return false, "no object selected", nil
	}
	opt := operationGet(ctx.Path())
	if opt == "" {
		return false, "not support operation", nil
	}

	var err error
	for _, v := range proxys {
		logger.Infof("start proxy %s operation %s", v.Tag, opt)
		switch opt {
			case OPERATION_START: {
				err = EnginStart(v.Tag)
			}
			case OPERATION_STOP: {
				err = EnginStop(v.Tag)
			}
			case OPERATION_RESTART: {
				EnginStop(v.Tag)
				err = EnginStart(v.Tag)
			}
		}

		if err != nil {
			logger.Errorf("proxy %s operation %s fail, %s", v.Tag, opt, err.Error())
			return false, "operation fail, " + err.Error(), nil
		}
	}

	return true, "operation success", nil
}

func ProxyTableGet(ctx *context.Context) (table.Table) {
	profile := table.NewDefaultTable(table.DefaultConfigWithDriver("sqlite"))

	profile.GetInfo().AddButton("Restart", icon.Refresh, action.Ajax(OPERATION_RESTART, enginProcess))
	profile.GetInfo().AddButton("Stop", icon.Pause, action.Ajax(OPERATION_STOP, enginProcess))
	profile.GetInfo().AddButton("Start", icon.Play, action.Ajax(OPERATION_START, enginProcess))

	info := profile.GetInfo().HideFilterArea().HideExportButton().HideFilterButton().HideRowSelector().HideQueryInfo()
	info.AddField("ID", "id", db.Int).FieldFilterable()
	info.AddField("Tag", "tag", db.Varchar).FieldFixed()

	info.AddField("Enable", "enable", db.Integer).FieldDisplay(func(model types.FieldModel) interface{} {
		return model.Value
	}).FieldEditAble(editType.Switch).FieldEditOptions(types.FieldOptions{
		{Value: "1", Text: "1"},
		{Value: "0", Text: "0"},
	}).FieldFilterable(types.FilterType{FormType: form.SelectSingle}).FieldFilterOptions(types.FieldOptions{
		{Value: "1", Text: "1"},
		{Value: "0", Text: "0"},
	})
	info.AddField("Auth", "auth", db.Integer).FieldDisplay(func(model types.FieldModel) interface{} {
		return model.Value
	}).FieldEditAble(editType.Switch).FieldEditOptions(types.FieldOptions{
		{Value: "1", Text: "1"},
		{Value: "0", Text: "0"},
	}).FieldFilterable(types.FilterType{FormType: form.SelectSingle}).FieldFilterOptions(types.FieldOptions{
		{Value: "1", Text: "1"},
		{Value: "0", Text: "0"},
	})

	info.AddField("Iface", "interface", db.Varchar)
	info.AddField("Port", "port", db.Integer)
	info.AddField("Protocal", "protocal", db.Varchar)
	info.AddField("Mode", "mode", db.Varchar)
	info.AddField("Remote", "remote", db.Varchar)

	info.AddField("Status", "status", db.Varchar).
		FieldDisplay(func(value types.FieldModel) interface{} {
			if value.Value == "" {
				return "stoped"
			}
			return value.Value
		}).
		FieldDot(map[string]types.FieldDotColor{
			"running": types.FieldDotColorInfo,
			"stoped": types.FieldDotColorDanger,
		}, types.FieldDotColorDanger)

	info.SetTable("proxys").SetTitle("Proxy Server Config").SetDescription("edit proxy servers config")

	addFrom := profile.GetForm()
	addFrom.AddField("Tag", "tag", db.Varchar, form.Text).
		FieldPostFilterFn(func(value types.PostFieldModel) interface{} {
			if value.Value.Value() == "" {
				logger.Error("user is null", value.Value.Value())
				return ""
			}
			return value.Value.Value()
		})

	var ifaceOptions types.FieldOptions = []types.FieldOption{
		{ Text: "0.0.0.0", Value: "0.0.0.0"},
	}
	ifaces, _ := util.InterfaceGet()
	for _, v := range ifaces {
		if v.IP == nil {
			continue
		}
		ifaceOptions = append(ifaceOptions, types.FieldOption{
			Text: v.IP.String(), Value: v.IP.String(),
		})
	}

	addFrom.AddField("Iface", "interface", db.Varchar, form.SelectSingle).
		FieldOptions(ifaceOptions).FieldRowWidth(2)

	addFrom.AddField("Port", "port", db.Tinyint, form.Number).
		FieldDefault("8080")

	var protcalOptions types.FieldOptions = []types.FieldOption {
		{Text: "HTTP", Value: "http",},
		{Text: "HTTPS", Value: "https",},
		{Text: "SOCK5", Value: "sock5",},
	}
	addFrom.AddField("Protocal", "protocal", db.Varchar, form.SelectSingle).
		FieldOptions(protcalOptions).FieldRowWidth(2)

	addFrom.AddField("Auth", "auth", db.Integer, form.Switch).
		FieldOptions(types.FieldOptions{
			{Value: "1", Text: "ON"},
			{Value: "0", Text: "OFF"},
		}).FieldDefault("0")

	var modeOptions types.FieldOptions = []types.FieldOption {
		{Text: "Local", Value: models.MODE_LOCAL,},
		{Text: "Remote", Value: models.MODE_REMOTE,},
		{Text: "Domain", Value: models.MODE_DOMAIN,},
	}
	addFrom.AddField("Mode", "mode", db.Varchar, form.SelectSingle).
		FieldOptions(modeOptions).FieldRowWidth(2)

	var remoteOptions types.FieldOptions
	remotes := models.RemoteGet()
	for _, v := range remotes {
		remoteOptions = append(remoteOptions, types.FieldOption{
			Text: v.Tag, Value: v.Tag,
		})
	}
	addFrom.AddField("Remote", "remote", db.Varchar, form.SelectSingle).
		FieldOptions(remoteOptions).FieldRowWidth(2)

	addFrom.AddField("Enable", "enable", db.Tinyint, form.Number).FieldDefault("1").FieldHide()
	addFrom.AddField("Status", "status", db.Varchar, form.Text).FieldDefault("stoped").FieldHide()

	addFrom.SetPostValidator(func(values form2.Values) error {
		if values.IsSingleUpdatePost() {
			proxy := models.ProxyFindByID(values.Get("id"))
			if proxy == nil {
				return fmt.Errorf("proxy id %s not exist", values.Get("id"))
			}
			if values.Has("enable") {
				return nil
			}
			if values.Has("auth") {
				EnginAuth(proxy.Tag, util.Atoi(values.Get("auth")))
				return nil
			}
			return fmt.Errorf("post single only enable / auth update")
		}

		if len(values.Get("tag")) < 3 {
			return fmt.Errorf("tag should more than 3 characters")
		}

		if values.Get("mode") != models.MODE_LOCAL {
			remote := values.Get("remote")
			if remote == "" && models.RemoteFind(remote) == nil {
				return fmt.Errorf("remote %s server config not exist", remote)
			}
		}

		port := util.Atoi(values.Get("port"))
		if port < 1 || port > 65535 {
			return fmt.Errorf("port %d is illegal", port)
		}

		go func() {
			tag := values.Get("tag")
			time.Sleep(time.Second)

			EnginStop(tag)
			EnginStart(tag)
		}()

		return nil
	})

	addFrom.SetTable("proxys").SetTitle("Proxy Server Config").SetDescription("edit proxy servers config")
	return profile
}

type MultiProxyCtrl struct {
	engin map[string] *ProxyEngin
	sync.RWMutex
}

var multiProxy MultiProxyCtrl

func EnginInit() {
	multiProxy.engin = make(map[string] *ProxyEngin, 100)
	proxys := models.ProxyGet()
	for _, v := range proxys {
		if 0 == v.Enable {
			continue
		}
		err := EnginStart(v.Tag)
		if err != nil {
			logger.Errorf("proxy %s start fail, %s", v.Tag, err.Error())
		} else {
			logger.Infof("proxy %s start success", v.Tag)
		}
	}

	go deleteSync()
	go statusSync()
}

func EnginFini()  {
	multiProxy.Lock()
	defer multiProxy.Unlock()

	for _, v := range multiProxy.engin {
		v.Stop()
	}
}

func EnginStart(tag string) error {
	multiProxy.Lock()
	defer multiProxy.Unlock()

	var err error

	engin, _ := multiProxy.engin[tag]
	if engin != nil {
		return fmt.Errorf("proxy %s engin running", tag)
	}

	proxy := models.ProxyFind(tag)
	if proxy == nil {
		return fmt.Errorf("proxy %s not exist", tag)
	}

	var remote *models.Remote
	if proxy.Mode != models.MODE_LOCAL {
		remote = models.RemoteFind(proxy.Remote)
		if remote == nil {
			errs := fmt.Sprintf("remote %s not exist", proxy.Remote)
			models.ProxyUpdate(proxy.Tag, func(u *models.Proxy) {
				u.Status = errs
			})
			return fmt.Errorf(errs)
		}
	}

	engin, err = NewProxyEngin(proxy, remote)
	if err != nil {
		models.ProxyUpdate(proxy.Tag, func(u *models.Proxy) {
			u.Status = err.Error()
		})
		return err
	}

	models.ProxyUpdate(proxy.Tag, func(u *models.Proxy) {
		u.Status = "running"
	})

	multiProxy.engin[tag] = engin
	return nil
}

func EnginAuth(tag string, auth int) error {
	multiProxy.Lock()
	defer multiProxy.Unlock()

	engin, _ := multiProxy.engin[tag]
	if engin == nil {
		return fmt.Errorf("proxy %s engin stoped", tag)
	}

	engin.AuthSwitch(auth)
	return nil
}

func EnginStop(tag string) error {
	multiProxy.Lock()
	defer multiProxy.Unlock()

	models.ProxyUpdate(tag, func(u *models.Proxy) {
		u.Status = "stoped"
	})

	engin, _ := multiProxy.engin[tag]
	if engin == nil {
		return fmt.Errorf("proxy %s engin stoped", tag)
	}

	engin.Stop()
	delete(multiProxy.engin, tag)
	return nil
}

func deleteSync()  {
	for  {
		time.Sleep(time.Second)

		multiProxy.Lock()
		for _, v := range multiProxy.engin {
			if models.ProxyFind(v.proxy.Tag) != nil {
				continue
			}
			logger.Warnf("engin %s has beed delete", v.proxy.Tag)
			v.Stop()
			delete(multiProxy.engin, v.proxy.Tag)
		}
		multiProxy.Unlock()
	}
}

func statusSync()  {
	for  {
		time.Sleep(time.Second)

		proxys := models.ProxyGet()
		for _, v := range proxys {
			if v.Enable == 0 {
				EnginStop(v.Tag)
			} else {
				EnginStart(v.Tag)
			}
		}
	}
}