package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/davecgh/go-spew/spew"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var (
	t   *template.Template
	err error
)

func unfilter(s string) string {
	return strings.Replace(s, "\n", `
	`, -1)
}

func main() {
	prerenderCount := 0

	renderBuffer := &bytes.Buffer{}

	templateFuncs := map[string]interface{}{
		"newl": func() string {
			return string('\n')
		},
		"tab": func() string {
			return string('\t')
		},
		"dump": spew.Sdump,
		"isMap": func(arg interface{}) bool {
			// TODO: if we'd want to check "is iterable" via reflection
			_, ok := arg.(map[string]interface{})
			return ok
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
			prerenderCount++

			var m map[string]interface{}
			switch len(args) {
			case 0:
			case 1:
				m = args[0]
			default:
				log.Fatal("Error: multiple arguments not supported, use map")
			}

			err = t.ExecuteTemplate(renderBuffer, templateName, m)
			check(err)

			result := renderBuffer.String()
			renderBuffer.Reset()
			return result
		},
		"regex": func(data, pattern string) []map[string]string {
			re := regexp.MustCompile(pattern)
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

	t, err = template.New("").Funcs(templateFuncs).ParseFiles("./typescript.tpl")
	check(err)

	err = t.ExecuteTemplate(os.Stdout, "Main", nil)
	check(err)

	// spew.Dump(regex.String())

	// //Find all exported functions.
	// regex.Reset()
	// err = t.ExecuteTemplate(regex, "ExportFunc", nil)
	// check(err)
	// re := regexp.MustCompile(regex.String())
	// println("Regex: ", regex.String())
	// matches := re.FindAllStringSubmatch(string(script), -1)
	// result := namedSubmatchMap(matches, re.SubexpNames())
	// spew.Dump(result)

	// //Find all exported classes.
	// regex.Reset()
	// err = t.ExecuteTemplate(regex, "ExportClass", nil)
	// check(err)
	// re = regexp.MustCompile(regex.String())
	// println("Regex: ", regex.String())
	// matches = re.FindAllStringSubmatch(string(script), -1)
	// result = namedSubmatchMap(matches, re.SubexpNames())
	// spew.Dump(result)

	spew.Dump("CALLED:", prerenderCount)
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
