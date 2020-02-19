package gin

type LogInfo struct {
	Path string
	MaxSize uint32
	MaxRolls uint32
}

var LogTemplate string = `
<seelog minlevel="debug" maxlevel="critical">
	<outputs formatid="main">
		<rollingfile type="size" filename="{{.Path}}" maxsize="{{.MaxSize}}" maxrolls="{{.MaxRolls}}"/>
		<console/>
	</outputs>
	<formats>
	    <format id="main" format="[%LEVEL][%Date(2006-01-02T15:04:05.000000)][%Project][%File:%Line][%FuncShort] => %Msg%n"/>
	</formats>
</seelog>
`
