package template

type CommonInfo struct {
	Project string
	Version string
}

var CommonTemplate string = `
package common

var ProjectName string = "{{.Project}}"
var ProjectVersion string = "{{.Version}}"
`
