package smallbox

import (
	"html/template"
	"strings"

	adminTemplate "github.com/GoAdminGroup/go-admin/template"
)

type SmallBox struct {
	*adminTemplate.BaseComponent

	Title      template.HTML
	Value      template.HTML
	Url        string
	Color      template.HTML
	IsSvg      bool
	IsHexColor bool
	Icon       template.HTML
}

func New() SmallBox {
	return SmallBox{
		BaseComponent: &adminTemplate.BaseComponent{
			Name:     "smallbox",
			HTMLData: List["smallbox"],
		},
	}
}

func (s SmallBox) SetTitle(value template.HTML) SmallBox {
	s.Title = value
	return s
}

func (s SmallBox) SetValue(value template.HTML) SmallBox {
	s.Value = value
	return s
}

func (s SmallBox) SetColor(value template.HTML) SmallBox {
	s.Color = value
	if strings.Contains(string(value), "#") {
		s.IsHexColor = true
	}
	return s
}

func (s SmallBox) SetIcon(value template.HTML) SmallBox {
	s.Icon = value
	if strings.Contains(string(value), "svg") {
		s.IsSvg = true
	}
	return s
}

func (s SmallBox) SetUrl(value string) SmallBox {
	s.Url = value
	return s
}

func (s SmallBox) GetContent() template.HTML { return s.GetContentWithData(s) }
