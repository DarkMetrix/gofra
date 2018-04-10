package template

type NamingJsonInfo struct {
	Project string
	Addr string
}

var NamingJsonTemplate string = `
{
  "locations":
  {
    "{{.Project}}":"local|{{.Addr}}"
  }
}
`

