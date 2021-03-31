package html

import (
	"fmt"
	"html/template"
	"strings"
)

type M map[string]string

type Style M
type Attribute M

func (s Style) String() template.HTML {
	res := ""
	for k, v := range s {
		res += k + ":" + v + ";"
	}
	if res != "" {
		res = ` style="` + res + `"`
	}
	return template.HTML(res)
}

func (s Attribute) String() template.HTML {
	res := ""
	for k, v := range s {
		res += k + `="` + v + `" `
	}
	if res != "" {
		res = ` ` + res[:len(res)-1]
	}
	return template.HTML(res)
}

type Element struct {
	Tag       template.HTML
	Content   template.HTML
	Style     Style
	Attribute Attribute
}

func BaseEl() Element {
	return Element{Style: make(map[string]string), Attribute: make(map[string]string)}
}

func (b Element) SetTag(tag template.HTML) Element {
	b.Tag = tag
	return b
}

func (b Element) SetContent(content template.HTML) Element {
	b.Content = content
	return b
}

func (b Element) SetStyle(key, value string) Element {
	b.Style[key] = value
	return b
}

func (b Element) SetClass(class ...string) Element {
	if b.Attribute["class"] != "" {
		b.Attribute["class"] += " " + strings.Join(class, " ")
	} else {
		b.Attribute["class"] += strings.Join(class, " ")
	}
	return b
}

func (b Element) SetId(id string) Element {
	b.Attribute["id"] = id
	return b
}

func (b Element) SetAttr(key, value string) Element {
	b.Attribute[key] = value
	return b
}

func (b Element) SetStyleAndAttr(ms []M) Element {
	if len(ms) > 0 {
		for k, v := range ms[0] {
			b.Style[k] = v
		}
	}
	if len(ms) > 1 {
		for k, v := range ms[1] {
			b.Attribute[k] = v
		}
	}
	return b
}

func (b Element) Get() template.HTML {
	return template.HTML(fmt.Sprintf(`<%s%s%s>%s</%s>`, b.Tag, b.Style.String(), b.Attribute.String(), b.Content, b.Tag))
}

func BodyEl() Element {
	return BaseEl().SetTag("body")
}

func Body(content template.HTML, ms ...M) template.HTML {
	return BodyEl().SetContent(content).SetStyleAndAttr(ms).Get()
}

func DivEl() Element {
	return BaseEl().SetTag("div")
}

func Div(content template.HTML, ms ...M) template.HTML {
	return DivEl().SetContent(content).SetStyleAndAttr(ms).Get()
}

func SpanEl() Element {
	return BaseEl().SetTag("span")
}

func Span(content template.HTML, ms ...M) template.HTML {
	return SpanEl().SetContent(content).SetStyleAndAttr(ms).Get()
}

func IEl() Element {
	return BaseEl().SetTag("i")
}

func I(ms ...M) template.HTML {
	return IEl().SetStyleAndAttr(ms).Get()
}

func PEl() Element {
	return BaseEl().SetTag("p")
}

func P(content template.HTML, ms ...M) template.HTML {
	return PEl().SetContent(content).SetStyleAndAttr(ms).Get()
}

func ButtonEl() Element {
	return BaseEl().SetTag("button")
}

func Button(content template.HTML, ms ...M) template.HTML {
	return ButtonEl().SetContent(content).SetStyleAndAttr(ms).Get()
}

func FormEl() Element {
	return BaseEl().SetTag("form")
}

func Form(content template.HTML, ms ...M) template.HTML {
	return FormEl().SetContent(content).SetStyleAndAttr(ms).Get()
}

func AEl() Element {
	return BaseEl().SetTag("a")
}

func A(content template.HTML, ms ...M) template.HTML {
	return AEl().SetContent(content).SetStyleAndAttr(ms).Get()
}

func UlEl() Element {
	return BaseEl().SetTag("ul")
}

func Ul(content template.HTML, ms ...M) template.HTML {
	return UlEl().SetContent(content).SetStyleAndAttr(ms).Get()
}

func LiEl() Element {
	return BaseEl().SetTag("li")
}

func Li(content template.HTML, ms ...M) template.HTML {
	return LiEl().SetContent(content).SetStyleAndAttr(ms).Get()
}

func BEl() Element {
	return BaseEl().SetTag("b")
}

func B(content template.HTML, ms ...M) template.HTML {
	return BEl().SetContent(content).SetStyleAndAttr(ms).Get()
}

func Br(num ...int) template.HTML {
	c := template.HTML("<br>")
	if len(num) > 0 {
		c = ""
		for range num {
			c += "<br>"
		}
	}
	return c
}

func H1El() Element {
	return BaseEl().SetTag("h1")
}

func H1(content template.HTML, ms ...M) template.HTML {
	return H1El().SetContent(content).SetStyleAndAttr(ms).Get()
}

func H2El() Element {
	return BaseEl().SetTag("h2")
}

func H2(content template.HTML, ms ...M) template.HTML {
	return H2El().SetContent(content).SetStyleAndAttr(ms).Get()
}

func H3El() Element {
	return BaseEl().SetTag("h3")
}

func H3(content template.HTML, ms ...M) template.HTML {
	return H3El().SetContent(content).SetStyleAndAttr(ms).Get()
}

func H4El() Element {
	return BaseEl().SetTag("h4")
}

func H4(content template.HTML, ms ...M) template.HTML {
	return H4El().SetContent(content).SetStyleAndAttr(ms).Get()
}

func H5El() Element {
	return BaseEl().SetTag("h5")
}

func H5(content template.HTML, ms ...M) template.HTML {
	return H5El().SetContent(content).SetStyleAndAttr(ms).Get()
}

func H6El() Element {
	return BaseEl().SetTag("h6")
}

func H6(content template.HTML, ms ...M) template.HTML {
	return H6El().SetContent(content).SetStyleAndAttr(ms).Get()
}
