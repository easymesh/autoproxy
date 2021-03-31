package description

import (
	"html/template"

	adminTemplate "github.com/GoAdminGroup/go-admin/template"
)

type Description struct {
	*adminTemplate.BaseComponent

	Border  string
	Number  template.HTML
	Title   template.HTML
	Arrow   string
	Color   template.HTML
	Percent template.HTML
}

func New() Description {
	return Description{
		BaseComponent: &adminTemplate.BaseComponent{
			Name:     "description",
			HTMLData: List["description"],
		},
	}
}

func (c Description) SetNumber(value template.HTML) Description {
	c.Number = value
	return c
}

func (c Description) SetTitle(value template.HTML) Description {
	c.Title = value
	return c
}

func (c Description) SetArrow(value string) Description {
	c.Arrow = value
	return c
}

func (c Description) SetPercent(value template.HTML) Description {
	c.Percent = value
	return c
}

func (c Description) SetColor(value template.HTML) Description {
	c.Color = value
	return c
}

func (c Description) SetBorder(value string) Description {
	c.Border = value
	return c
}

func (c Description) GetContent() template.HTML { return c.GetContentWithData(c) }
