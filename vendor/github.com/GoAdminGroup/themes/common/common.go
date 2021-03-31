package common

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/GoAdminGroup/go-admin/modules/config"
	adminTemplate "github.com/GoAdminGroup/go-admin/template"
)

type BaseTheme struct {
	AssetPaths   map[string]string
	TemplateList map[string]string
	Separation   bool
}

const Version = "v0.0.42"

func (b *BaseTheme) GetVersion() string {
	return Version
}

func (b *BaseTheme) GetRequirements() []string {
	return []string{">=v1.2.19"}
}

var comps = []string{"datatable.min.js", "form.min.js", "tree.min.js", "treeview.min.js"}

func inArray(s string, arr []string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}

func (b *BaseTheme) GetAssetImportHTML(exclude ...string) template.HTML {
	res := template.HTML("")
	for name, html := range b.AssetPaths {
		if inArray(name, comps) {
			if !inArray(strings.Replace(name, ".min.js", "", -1), exclude) {
				res += GetImportJSTag("/assets" + html)
			}
		}
	}
	return res
}

func (b *BaseTheme) GetHeadHTML() template.HTML {
	res := GetImportJSTag("/assets" + b.AssetPaths["all.min.js"])
	res += GetImportCSSTag("/assets" + b.AssetPaths["all.min.css"])
	return res
}

func (b *BaseTheme) GetFootJS() template.HTML {
	return GetImportJSTag("/assets" + b.AssetPaths["all_2.min.js"])
}

func (b *BaseTheme) getHTMLFromTmplList(key string) template.HTML {
	if !b.Separation {
		return template.HTML(b.TemplateList[key])
	}
	t, _ := template.ParseFiles(config.GetAssetRootPath() + "pages/" + b.TemplateList[key] + ".tmpl")
	var buf = new(bytes.Buffer)
	_ = t.Execute(buf, nil)
	return template.HTML(buf.String())
}

func (b *BaseTheme) Get500HTML() template.HTML {
	return b.getHTMLFromTmplList("500")
}

func (b *BaseTheme) Get404HTML() template.HTML {
	return b.getHTMLFromTmplList("400")
}

func (b *BaseTheme) Get403HTML() template.HTML {
	return b.getHTMLFromTmplList("403")
}

func (b *BaseTheme) GetTemplate(isPjax bool) (tmpl *template.Template, name string) {
	var err error

	if b.Separation {
		root := config.GetAssetRootPath() + "pages/"
		if !isPjax {
			name = "layout"
			tmpl, err = template.New("layout").Funcs(adminTemplate.DefaultFuncMap).
				ParseFiles(
					root+b.TemplateList["layout"]+".tmpl",
					root+b.TemplateList["head"]+".tmpl",
					root+b.TemplateList["header"]+".tmpl",
					root+b.TemplateList["sidebar"]+".tmpl",
					root+b.TemplateList["footer"]+".tmpl",
					root+b.TemplateList["js"]+".tmpl",
					root+b.TemplateList["menu"]+".tmpl",
					root+b.TemplateList["admin_panel"]+".tmpl",
					root+b.TemplateList["content"]+".tmpl")
		} else {
			name = "content"
			tmpl, err = template.New("content").Funcs(adminTemplate.DefaultFuncMap).
				ParseFiles(root+b.TemplateList["admin_panel"]+".tmpl", root+b.TemplateList["content"]+".tmpl")
		}
	} else {
		if !isPjax {
			name = "layout"
			tmpl, err = template.New("layout").Funcs(adminTemplate.DefaultFuncMap).
				Parse(b.TemplateList["layout"] +
					b.TemplateList["head"] + b.TemplateList["header"] + b.TemplateList["sidebar"] +
					b.TemplateList["footer"] + b.TemplateList["js"] + b.TemplateList["menu"] +
					b.TemplateList["admin_panel"] + b.TemplateList["content"])
		} else {
			name = "content"
			tmpl, err = template.New("content").Funcs(adminTemplate.DefaultFuncMap).
				Parse(b.TemplateList["admin_panel"] + b.TemplateList["content"])
		}
	}

	if err != nil {
		panic(err)
	}

	return
}

func GetImportJSTag(src string) template.HTML {
	if config.GetAssetUrl() != "" {
		return template.HTML(`<script src="` + config.GetAssetUrl() + src + `"></script>`)
	} else {
		return template.HTML(`<script src="` + config.Url(src) + `"></script>`)
	}
}

func GetImportCSSTag(src string) template.HTML {
	if config.GetAssetUrl() != "" {
		return template.HTML(`<link rel="stylesheet" href="` + config.GetAssetUrl() + src + `">`)
	} else {
		return template.HTML(`<link rel="stylesheet" href="` + config.Url(src) + `">`)
	}
}

var SepTemplateList = map[string]string{
	"403":                               "403",
	"404":                               "404",
	"500":                               "500",
	"admin_panel":                       "admin_panel",
	"components/alert":                  "components/alert",
	"components/box":                    "components/box",
	"components/button":                 "components/button",
	"components/col":                    "components/col",
	"components/form/array":             "components/form/array",
	"components/form/checkbox":          "components/form/checkbox",
	"components/form/checkbox_single":   "components/form/checkbox_single",
	"components/form/checkbox_stacked":  "components/form/checkbox_stacked",
	"components/form/code":              "components/form/code",
	"components/form/color":             "components/form/color",
	"components/form/currency":          "components/form/currency",
	"components/form/custom":            "components/form/custom",
	"components/form/datetime":          "components/form/datetime",
	"components/form/datetime_range":    "components/form/datetime_range",
	"components/form/default":           "components/form/default",
	"components/form/email":             "components/form/email",
	"components/form/file":              "components/form/file",
	"components/form/help_block":        "components/form/help_block",
	"components/form/iconpicker":        "components/form/iconpicker",
	"components/form/ip":                "components/form/ip",
	"components/form/multi_file":        "components/form/multi_file",
	"components/form/number":            "components/form/number",
	"components/form/number_range":      "components/form/number_range",
	"components/form/password":          "components/form/password",
	"components/form/radio":             "components/form/radio",
	"components/form/rate":              "components/form/rate",
	"components/form/richtext":          "components/form/richtext",
	"components/form/select":            "components/form/select",
	"components/form/selectbox":         "components/form/selectbox",
	"components/form/singleselect":      "components/form/singleselect",
	"components/form/slider":            "components/form/slider",
	"components/form/switch":            "components/form/switch",
	"components/form/table":             "components/form/table",
	"components/form/text":              "components/form/text",
	"components/form/textarea":          "components/form/textarea",
	"components/form/url":               "components/form/url",
	"components/form":                   "components/form",
	"components/form_components":        "components/form_components",
	"components/form_components_layout": "components/form_components_layout",
	"components/form_layout_default":    "components/form_layout_default",
	"components/form_layout_flow":       "components/form_layout_flow",
	"components/form_layout_tab":        "components/form_layout_tab",
	"components/form_layout_two_col":    "components/form_layout_two_col",
	"components/image":                  "components/image",
	"components/label":                  "components/label",
	"components/link":                   "components/link",
	"components/paginator":              "components/paginator",
	"components/popup":                  "components/popup",
	"components/row":                    "components/row",
	"components/table/box-header":       "components/table/box-header",
	"components/table":                  "components/table",
	"components/tabs":                   "components/tabs",
	"components/tree-header":            "components/tree-header",
	"components/tree":                   "components/tree",
	"components/treeview":               "components/treeview",
	"content":                           "content",
	"control_panel":                     "control_panel",
	"footer":                            "footer",
	"head":                              "head",
	"header":                            "header",
	"js":                                "js",
	"layout":                            "layout",
	"menu":                              "menu",
	"sidebar":                           "sidebar",
}
