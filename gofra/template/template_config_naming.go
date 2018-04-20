package template

type NamingTomlInfo struct {
	Project string
	Addr string
}

var NamingTomlTemplate string = `
# Naming configuration
#
# locations save all the setting of naming mapping
# You could implement your own version, by default a 'local' naming type is support
#
# Format:
#	NAMING_TYPE|NAMING_STRING
# eg:
#	local|127.0.0.1:58888
[locations]
	{{.Project}}="local|{{.Addr}}"
`

