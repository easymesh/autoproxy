package infobox

import (
	"html/template"
	"strings"

	adminTemplate "github.com/GoAdminGroup/go-admin/template"
)

type InfoBox struct {
	*adminTemplate.BaseComponent

	Icon       template.HTML
	Text       template.HTML
	Number     template.HTML
	Content    template.HTML
	Color      template.HTML
	IsHexColor bool
	IsSvg      bool
}

func New() InfoBox {
	return InfoBox{
		BaseComponent: &adminTemplate.BaseComponent{
			Name:     "infobox",
			HTMLData: List["infobox"],
		},
	}
}

func (i InfoBox) SetIcon(value template.HTML) InfoBox {
	i.Icon = value
	if strings.Contains(string(value), "svg") {
		i.IsSvg = true
	}
	return i
}

func (i InfoBox) SetText(value template.HTML) InfoBox {
	i.Text = value
	return i
}

func (i InfoBox) SetNumber(value template.HTML) InfoBox {
	i.Number = value
	return i
}

func (i InfoBox) SetContent(value template.HTML) InfoBox {
	i.Content = value
	return i
}

func (i InfoBox) SetColor(value template.HTML) InfoBox {
	i.Color = value
	if strings.Contains(string(value), "#") {
		i.IsHexColor = true
	}
	return i
}

func (i InfoBox) GetContent() template.HTML { return i.GetContentWithData(i) }
