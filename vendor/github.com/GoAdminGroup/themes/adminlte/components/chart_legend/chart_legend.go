package chart_legend

import (
	"html/template"

	adminTemplate "github.com/GoAdminGroup/go-admin/template"
)

type ChartLegend struct {
	*adminTemplate.BaseComponent

	Data []map[string]string
}

func New() ChartLegend {
	return ChartLegend{
		BaseComponent: &adminTemplate.BaseComponent{
			Name:     "chart-legend",
			HTMLData: List["chart-legend"],
		},
	}
}

func (c ChartLegend) SetData(value []map[string]string) ChartLegend {
	c.Data = value
	return c
}

func (c ChartLegend) GetContent() template.HTML { return c.GetContentWithData(c) }
