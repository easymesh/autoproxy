package pages

import (
	"fmt"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/logger"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	editType "github.com/GoAdminGroup/go-admin/template/types/table"
	"github.com/easymesh/autoproxy/console/models"
	"github.com/easymesh/autoproxy/console/uitl"
)

func ProxyTableGet(ctx *context.Context) (table.Table) {
	profile := table.NewDefaultTable(table.DefaultConfigWithDriver("sqlite"))

	info := profile.GetInfo().HideFilterArea().HideExportButton().HideFilterButton().HideRowSelector().HideQueryInfo()
	info.AddField("ID", "id", db.Int).FieldFilterable()
	info.AddField("Tag", "tag", db.Varchar)

	info.AddField("Enable", "enable", db.Integer).FieldDisplay(func(model types.FieldModel) interface{} {
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
	info.AddField("Protocal", "protocal", db.Varchar).FieldFilterable()

	info.AddField("Auth", "auth", db.Integer).FieldDisplay(func(model types.FieldModel) interface{} {
		return model.Value
	}).FieldEditAble(editType.Switch).FieldEditOptions(types.FieldOptions{
		{Value: "1", Text: "1"},
		{Value: "0", Text: "0"},
	}).FieldFilterable(types.FilterType{FormType: form.SelectSingle}).FieldFilterOptions(types.FieldOptions{
		{Value: "1", Text: "1"},
		{Value: "0", Text: "0"},
	})

	info.AddField("Status", "status", db.Varchar).
		FieldDisplay(func(value types.FieldModel) interface{} {
			if value.Value == "" {
				return "unkown"
			}
			return value.Value
		}).
		FieldDot(map[string]types.FieldDotColor{
			"connected": types.FieldDotColorInfo,
			"unkown": types.FieldDotColorDanger,
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
		{Text: "Local", Value: "local",},
		{Text: "Remote", Value: "remote",},
		{Text: "Domain", Value: "domain",},
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
	addFrom.SetPostValidator(func(values form2.Values) error {
		if values.IsSingleUpdatePost() {
			if !values.Has("enable") {
				return fmt.Errorf("account single only enable update")
			}
			return nil
		}
		if len(values.Get("tag")) < 3 {
			return fmt.Errorf("tag should more than 3 characters")
		}
		remote := values.Get("remote")
		if remote == "" && models.RemoteFind(remote) == nil {
			return fmt.Errorf("remote server config not exist", remote)
		}
		port := util.Atoi(values.Get("port"))
		if port < 1 || port > 65535 {
			return fmt.Errorf("port %d is illegal", port)
		}


		//if len(values.Get("address")) < 3 {
		//	return fmt.Errorf("tag should more than 3 characters")
		//}
		//if len(values.Get("domains")) < 1 {
		//	return fmt.Errorf("domains should more than 1 characters")
		//}
		return nil
	})



	addFrom.SetTable("proxys").SetTitle("Proxy Server Config").SetDescription("edit proxy servers config")
	return profile
}