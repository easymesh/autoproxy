package chart_legend

var List = map[string]string{
	"chart-legend": `{{define "chart-legend"}}
<ul class="chart-legend clearfix">
    {{range $key, $data := .Data}}
        <li><i class="fa fa-circle-o text-{{index $data "color"}}"></i>{{index $data "label"}}</li>
    {{end}}
</ul>
{{end}}`,
}
