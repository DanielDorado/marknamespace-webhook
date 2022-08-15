package main

import (
	"bytes"
	"html/template"
	"regexp"

	"k8s.io/klog/v2"
)

func processNamespace(namespace string, rules []Rule) map[string]string {
	labels := map[string]string{}

	for _, r := range rules {
		re, err := regexp.Compile(r.CaseNamespace)
		if err != nil {
			klog.Error("Processing rules. Case %s is a wrong Regular Expression: %e", r.CaseNamespace, err)
			continue
		}
		values := re.FindStringSubmatch(namespace) // return [string that matches, submatch 0, submatch 1, ...]
		if values == nil || len(values) == 0 || len(values) == 1 && values[0] == "" {
			continue
		}
		for _, l := range r.Inject {
			templateValues := values[1:]
			name, err := fill(l.Name, templateValues)
			if err != nil {
				klog.Error("Processing template name: %s: %e", l.Name, err)
				continue
			}
			value, err := fill(l.Value, templateValues)
			if err != nil {
				klog.Error("Processing template value: %s: %e", l.Value, err)
				continue
			}
			labels[name] = value
		}

	}
	return labels
}

func fill(tpl string, values interface{}) (string, error) {
	tmpl, err := template.New("test").Parse(tpl)
	if err != nil {
		klog.Errorf("Processing template tpl: %s with values: %s", tpl, values)
		return "", err
	}
	var doc bytes.Buffer
	err = tmpl.Execute(&doc, values)
	if err != nil {
		klog.Errorf("Processing template tpl: %s with values: %s", tpl, values)
		return "", err
	}

	return doc.String(), nil
}
