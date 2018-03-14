package template

type CommonInfo struct {
	Author string
	Time string
	Project string
	Version string
}

var CommonTemplate string = `
/**********************************
 * Author : {{.Author}}
 * Time : {{.Time}}
 **********************************/

package common

var ProjectName string = "{{.Project}}"
var ProjectVersion string = "{{.Version}}"
`
