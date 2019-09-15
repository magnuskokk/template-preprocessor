package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"text/template"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	var t *template.Template
	regexCache := make(map[string]*regexp.Regexp)
	templateFuncs := map[string]interface{}{
		"newl": func() string {
			return string('\n')
		},
		"tab": func() string {
			return string('\t')
		},
		"map": func(args ...interface{}) map[string]interface{} {
			length := len(args)
			if length%2 != 0 {
				panic("Must have even number of arguments as key value pairs: (map \"arg1\" \"val1\" \"arg2\" \"val2\")")
			}

			m := make(map[string]interface{}, length/2)
			for i, j := 0, 1; i < length; i, j = i+2, j+2 {
				key := args[i].(string)
				val := args[j].(string)
				m[key] = val
			}
			return m
		},
		"render": func(templateName string, args ...map[string]interface{}) string {
			var m map[string]interface{}
			switch len(args) {
			case 0:
			case 1:
				m = args[0]
			default:
				log.Fatal("Error: multiple arguments not supported, use map")
			}

			b := &bytes.Buffer{}
			err := t.ExecuteTemplate(b, templateName, m)
			check(err)
			return b.String()
		},
		"regex": func(data, pattern string) []map[string]string {
			var re *regexp.Regexp
			var ok bool
			if re, ok = regexCache[pattern]; !ok {
				re = regexp.MustCompile(pattern)
				regexCache[pattern] = re
			}

			matches := re.FindAllStringSubmatch(data, -1)
			result := namedSubmatchMaps(matches, re.SubexpNames())
			return result
		},
		"loadFile": func(filename string) string {
			content, err := ioutil.ReadFile(filename)
			check(err)
			return string(content)
		},
	}

	var err error
	t, err = template.New("").Funcs(templateFuncs).ParseFiles("./main.tmpl")
	check(err)

	err = t.ExecuteTemplate(os.Stdout, "Main", nil)
	check(err)
}

func namedSubmatchMaps(matches [][]string, subexpNames []string) []map[string]string {
	var result []map[string]string
	for _, match := range matches {
		m := make(map[string]string)
		for i, name := range subexpNames {
			if i != 0 && name != "" {
				m[name] = string(match[i])
			}
		}
		result = append(result, m)
	}
	return result
}
