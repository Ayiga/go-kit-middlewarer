package main

import (
	"fmt"
	"text/template"
)

var templateFuncs template.FuncMap = map[string]interface{}{
	"importFilter": func(imp string, toFilterOut ...string) string {
		fmt.Println("Received arguments: ", imp, " and ", toFilterOut)
		for _, s := range toFilterOut {
			if imp == s {
				fmt.Println("Filtering out: ", s)
				return ""
			}
		}

		fmt.Println("Keeping: ", imp)
		return imp
	},
}
