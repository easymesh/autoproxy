package progress_group

import (
	"html/template"
	"strings"

	adminTemplate "github.com/GoAdminGroup/go-admin/template"
)

type ProgressGroup struct {
	*adminTemplate.BaseComponent

	Title       template.HTML
	Molecular   int
	Denominator int
	Color       template.HTML
	IsHexColor  bool
	Percent     int
}

func New() ProgressGroup {
	return ProgressGroup{
		BaseComponent: &adminTemplate.BaseComponent{
			Name:     "progress-group",
			HTMLData: List["progress-group"],
		},
	}
}

func (p ProgressGroup) SetTitle(value template.HTML) ProgressGroup {
	p.Title = value
	return p
}

func (p ProgressGroup) SetColor(value template.HTML) ProgressGroup {
	p.Color = value
	if strings.Contains(string(value), "#") {
		p.IsHexColor = true
	}
	return p
}

func (p ProgressGroup) SetPercent(value int) ProgressGroup {
	p.Percent = value
	return p
}

func (p ProgressGroup) SetDenominator(value int) ProgressGroup {
	p.Denominator = value
	return p
}

func (p ProgressGroup) SetMolecular(value int) ProgressGroup {
	p.Molecular = value
	return p
}

func (p ProgressGroup) GetContent() template.HTML { return p.GetContentWithData(p) }
