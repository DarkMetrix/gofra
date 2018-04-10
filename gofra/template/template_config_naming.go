package template

type NamingJsonInfo struct {
	Project string
	Addr string
}

var NamingJsonTemplate string = `
{
  "locations":
  {
    "{{.Project}}":
	{
      "is_test":false,
      "location_test":"{{.Addr}}",
      "location_real":"{{.Addr}}"
	}
  }
}
`

