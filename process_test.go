package main

import (
	"fmt"
	"testing"
)

func TestGetLabelFromNamespace(t *testing.T) {

	labelsConfiguration := `
labels:
  - caseNamespace: "^([^-]+)-([^-]+)-([^-]+)-([^-]+)$"
    inject:
    - name: "region"
      value: "reg-{{index . 0}}"
    - name: area
      value: "{{index . 1}}"
    - name: team
      value: "{{index . 2}}"
    - name: environment
      value: "{{index . 3}}"
  - caseNamespace: "^([^-]+)-([^-]+)-([^-]+)$"
    inject:
    - name: region
      value: "reg-{{index . 0}}"
    - name: area
      value: "{{index . 1}}"
    - name: team
      value: "{{index . 2}}"
    - name: environment
      value: "prod"
  - caseNamespace: "^([^-]+)-([^-]+)$"
    inject:
    - name: region
      value: "reg-{{index . 0}}"
    - name: system
      value: "sys-{{index . 1}}"
`
	config, err := NewConfig([]byte(labelsConfiguration))
	if err != nil {
		fmt.Printf("ERROR: %e", err)
	}
	for _, i := range []struct {
		Namespace      string
		ExpectedLabels string
	}{
		{
			"region1-area1-team1-environment1",
			"map[area:area1 environment:environment1 region:reg-region1 team:team1]",
		},
		{
			"area1-team1-environment1",
			"map[area:team1 environment:prod region:reg-area1 team:environment1]",
		},
		{
			"region1-system1",
			"map[region:reg-region1 system:sys-system1]",
		},
	} {
		gotLabels := processNamespace(i.Namespace, config.Labels)
		if fmt.Sprintf("%+v", gotLabels) != i.ExpectedLabels {
			t.Errorf("\nexpected %#v, got %#v", i.ExpectedLabels, fmt.Sprintf("%v", gotLabels))
		}
	}
}
