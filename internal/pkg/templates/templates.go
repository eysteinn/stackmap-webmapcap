package templates

// "wms_timeextent "{{range .Times}}{{.}},{{end}}{{.StartRange}}/{{.EndRange}}"
const timesTemplate = `
	
	"wms_timeextent "{{.AllTimes}},{{.StartRange}}/{{.EndRange}}
	"wms_timedefault" "{{.TimeDefault}}"
`
