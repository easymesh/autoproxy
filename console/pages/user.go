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

func UserTableGet(ctx *context.Context) (table.Table) {
	profile := table.NewDefaultTable(table.DefaultConfigWithDriver("sqlite"))

	info := profile.GetInfo().HideFilterArea().HideExportButton().HideFilterButton().HideRowSelector().HideQueryInfo()
	info.AddField("ID", "id", db.Int).FieldFilterable()
	info.AddField("User", "user", db.Varchar)
	info.AddField("Password", "password", db.Varchar)
	info.AddField("Enable", "enable", db.Integer).FieldDisplay(func(model types.FieldModel) interface{} {
		return model.Value
	}).FieldEditAble(editType.Switch).FieldEditOptions(types.FieldOptions{
		{Value: "1", Text: "1"},
		{Value: "0", Text: "0"},
	}).FieldFilterable(types.FilterType{FormType: form.SelectSingle}).FieldFilterOptions(types.FieldOptions{
		{Value: "1", Text: "1"},
		{Value: "0", Text: "0"},
	})
	//info.AddField("Flow", "flow", db.Integer).FieldDisplay(func(value types.FieldModel) interface{} {
	//	return fmt.Sprintf("%sGB", value.Value)
	//})
	//info.AddField("Online", "online", db.Varchar).
	//	FieldDisplay(func(value types.FieldModel) interface{} {
	//		if value.Value == "1" {
	//			return "online"
	//		}
	//		return "offline"
	//	}).
	//	FieldDot(map[string]types.FieldDotColor{
	//		"offline": types.FieldDotColorDanger,
	//		"online": types.FieldDotColorInfo,
	//	}, types.FieldDotColorDanger)

	info.SetTable("users").SetTitle("User").SetDescription("edit user account")

	addFrom := profile.GetForm()
	addFrom.AddField("User", "user", db.Varchar, form.Text).
		FieldPostFilterFn(func(value types.PostFieldModel) interface{} {
			if value.Value.Value() == "" {
				logger.Error("name is null", value.Value.Value())
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
	addFrom.AddField("Flow", "flow", db.Tinyint, form.Number).FieldDefault("0").FieldHide()
	addFrom.AddField("Online", "online", db.Tinyint, form.Number).FieldDefault("0").FieldHide()

	addFrom.SetPostValidator(func(values form2.Values) error {
		if values.IsSingleUpdatePost() {
			if !values.Has("enable") {
				return fmt.Errorf("account single only enable update")
			}
			return nil
		}
		if len(values.Get("user")) < 5 {
			return fmt.Errorf("account user should more than 5 characters")
		}
		if len(values.Get("password")) < 8 {
			return fmt.Errorf("account password should more than 8 characters")
		}
		return nil
	})
	addFrom.SetTable("users").SetTitle("User").SetDescription("edit user account")
	return profile
}

