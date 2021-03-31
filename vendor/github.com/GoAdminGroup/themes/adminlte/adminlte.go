package adminlte

import (
	"strings"

	adminTemplate "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/components"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/themes/adminlte/resource"
	"github.com/GoAdminGroup/themes/common"
	"github.com/gobuffalo/packr/v2"
)

const (
	ColorschemeSkinBlack       = "skin-black"
	ColorschemeSkinBlackLight  = "skin-black-light"
	ColorschemeSkinBlue        = "skin-blue"
	ColorschemeSkinBlueLight   = "skin-blue-light"
	ColorschemeSkinGreen       = "skin-green"
	ColorschemeSkinGreenLight  = "skin-green-light"
	ColorschemeSkinPurple      = "skin-purple"
	ColorschemeSkinPurpleLight = "skin-purple-light"
	ColorschemeSkinRed         = "skin-red"
	ColorschemeSkinRedLight    = "skin-red-light"
	ColorschemeSkinYellow      = "skin-yellow"
	ColorschemeSkinYellowLight = "skin-yellow-light"
)

type Theme struct {
	ThemeName string
	components.Base
	*common.BaseTheme
}

var Adminlte = Theme{
	ThemeName: "adminlte",
	Base: components.Base{
		Attribute: types.Attribute{
			TemplateList: TemplateList,
		},
	},
	BaseTheme: &common.BaseTheme{
		AssetPaths:   resource.AssetPaths,
		TemplateList: TemplateList,
	},
}

func init() {
	adminTemplate.Add("adminlte", &Adminlte)
}

func Get() *Theme {
	return &Adminlte
}

func (t *Theme) Name() string {
	return t.ThemeName
}

func (t *Theme) GetTmplList() map[string]string {
	return TemplateList
}

func (t *Theme) GetAsset(path string) ([]byte, error) {
	path = strings.Replace(path, "/assets/dist", "", -1)
	box := packr.New("adminlte", "./resource/assets/dist")
	return box.Find(path)
}

func (t *Theme) GetAssetList() []string {
	return resource.AssetsList
}
