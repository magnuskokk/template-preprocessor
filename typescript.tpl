
{{define "TestTemplate"}}
{{ render "Capture" (map "Name" "ident" "Pattern" "patternVal") }}
{{end}}


{{define "SetVars"}}
	{{range $key, $value := .}}
	{{`{{$`}}{{$key}} := "{{$value}}"{{`}}`}}
	{{end}}
{{end}}

{{define "UnrollVars"}}
	{{- range $key, $value := .}}"{{$key}}" "{{$value}}" {{end -}}
{{end}}

{{define "Capture"}}(?P<{{.Name}}>{{.Pattern}}){{end}}

{{define "RecTemplate"}}
{{template "TestTemplate"}}
{{end}}

{{define "IdentPattern"}}.*?{{end}}

{{define "FuncRegex"}}
	{{- $ident := render "IdentPattern"}}
	{{- render "Capture" (map "Name" "FuncName" "Pattern" $ident) }}\(
		{{- $ident -}}
	\):\s
	{{- render "Capture" (map "Name" "FuncReturn" "Pattern" $ident) }};
{{end}}
{{define "ExportFunc"}}export function {{template "FuncRegex"}}{{end}}

{{define "ClassRegex"}}
	{{- $ident := render "IdentPattern"}}
	{{- $name := map "Name" "ClassName" "Pattern" $ident}}
	{{- $body := map "Name" "ClassBody" "Pattern" ".*?[^}]*"}}
	{{- render "Capture" $name}}\s{
		{{- render "Capture" $body -}}
	}
{{end}}
{{define "ExportClassRegex"}}export class {{ template "ClassRegex"}}{{end}}

{{/* Main parses the typescript file and prints classes and functions */}}
{{define "Main"}}
{{- $classRegex := render "ClassRegex"}}
{{- $funcRegex := render "FuncRegex"}}
{{- $file := loadFile "./test.d.ts"}}

{{- range $classMap := regex $file $classRegex}}
	{{- print "Class name: "}}{{$classMap.ClassName}}

	{{- print "Functions: "}}
	{{- range $funcMap := regex $classMap.ClassBody $funcRegex}}
		{{- newl}}
		{{- tab}}{{ print "Name: "}}{{$funcMap.FuncName}}
		{{- newl}}
		{{- tab}}{{ print "Return: "}}{{$funcMap.FuncReturn}}
	{{end}}
{{end}}
{{end}}
