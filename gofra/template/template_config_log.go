package template

type LogInfo struct {
	Path string
	MaxSize uint32
	MaxRolls uint32
}

var LogTemplate string = `
<seelog>
	<outputs formatid="main">
		<rollingfile type="size" filename="{{.Path}}" maxsize="{{.MaxSize}}" maxrolls="{{.MaxRolls}}"/>
	</outputs>
	<formats>
	    <format id="main" format="[%LEVEL][%DateT%Time][%File:%Line][%FuncShort] => %Msg%n"/>
	</formats>
</seelog>
`
