package jdenticon

const tmpl = `<svg width="{{.Width}}" height="{{.Height}}" preserveAspectRatio="xMidYMid meet" viewBox="0 0 {{.Width}} {{.Height}}" xmlns="http://www.w3.org/2000/svg">
	{{- range .Paths -}}
		<path 
			{{- if .Fill}} fill="{{.Fill}}"{{end -}}
			{{- if .Stroke}} stroke="{{.Stroke}}"{{end -}}
			{{- if .UseOpacity}} opacity="{{.Opacity}}"{{end -}}
			{{- if .Shapes}} d="{{.Shapes}}"{{end -}}
		/>
	{{- end -}}
</svg>`
