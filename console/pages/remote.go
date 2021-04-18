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
)

func RemoteTableGet(ctx *context.Context) (table.Table) {
	profile := table.NewDefaultTable(table.DefaultConfigWithDriver("sqlite"))

	info := profile.GetInfo().HideFilterArea().HideExportButton().HideFilterButton().HideRowSelector().HideQueryInfo()
	info.AddField("ID", "id", db.Int).FieldFilterable()
	info.AddField("Address", "address", db.Varchar).FieldFilterable()
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
	info.AddField("User", "user", db.Varchar).FieldFilterable()
	info.AddField("Password", "password", db.Varchar).FieldFilterable()
	info.AddField("Protocal", "protocal", db.Varchar).FieldFilterable()
	info.SetTable("remotes").SetTitle("Remote Server Config").SetDescription("edit remote servers config")

	addFrom := profile.GetForm()
	addFrom.AddField("Address", "address", db.Varchar, form.Text).
		FieldPostFilterFn(func(value types.PostFieldModel) interface{} {
			if value.Value.Value() == "" {
				logger.Error("address is null", value.Value.Value())
				return ""
			}
			return value.Value.Value()
		})

	var protcalOptions types.FieldOptions
	protcalOptions = append(protcalOptions, types.FieldOption{
		Text: "HTTP", Value: "http",
	})
	protcalOptions = append(protcalOptions, types.FieldOption{
		Text: "HTTPS", Value: "https",
	})
	protcalOptions = append(protcalOptions, types.FieldOption{
		Text: "SOCK5", Value: "sock5",
	})

	addFrom.AddField("Protocal", "protocal", db.Varchar, form.SelectSingle).
		FieldOptions(protcalOptions).FieldRowWidth(2)

	addFrom.AddField("Auth", "auth", db.Integer, form.Switch).
		FieldOptions(types.FieldOptions{
			{Value: "1", Text: "ON"},
			{Value: "0", Text: "OFF"},
		}).FieldDefault("1")

	addFrom.AddField("User", "user", db.Varchar, form.Text).
		FieldPostFilterFn(func(value types.PostFieldModel) interface{} {
			if value.Value.Value() == "" {
				logger.Error("user is null", value.Value.Value())
				return ""
			}
			return value.Value.Value()
		})

	addFrom.AddField("Password", "password", db.Varchar, form.Password).
		FieldPostFilterFn(func(value types.PostFieldModel) interface{} {
			if value.Value.Value() == "" {
				logger.Error("password is null", value.Value.Value())
				return ""
			}
			return value.Value.Value()
		}).FieldDefault("").FieldValue("")

	addFrom.AddField("Enable", "enable", db.Tinyint, form.Number).FieldDefault("1").FieldHide()
	addFrom.SetPostValidator(func(values form2.Values) error {
		if values.IsSingleUpdatePost() {
			if !values.Has("enable") {
				return fmt.Errorf("account single only enable update")
			}
			return nil
		}
		//if len(values.Get("address")) < 3 {
		//	return fmt.Errorf("tag should more than 3 characters")
		//}
		//if len(values.Get("domains")) < 1 {
		//	return fmt.Errorf("domains should more than 1 characters")
		//}
		return nil
	})
	addFrom.SetTable("remotes").SetTitle("Remote Server Config").SetDescription("edit remote servers config")
	return profile
}