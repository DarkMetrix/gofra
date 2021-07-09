package option

// Options represents the option information
type Options struct {
	OutputPath string

	// general & docker options
	Project string
	Author  string

	// override control
	Override    bool
	IgnoreExist bool

	// go information
	GoModule  string
	GoVersion string

	// protobuf information
	ProtocPath           string
	ProtoFileIncludePath []string

	// service information
	Addr                string
	PackageName         string
	ServiceName         string
	RPCName             string
	RequestName         string
	ResponseName        string
	ImportedPackageName string
	ServiceImplName     string

	// config package path
	ConfigPackagePath string

	// istio & kubernetes
	Namespace  string
	Version    string
	Port       string
	TargetPort string
	ImagePath  string
}

// NewOptions returns a new Options pointer
func NewOptions(opts ...Option) *Options {
	options := &Options{}
	for _, optionFunc := range opts {
		optionFunc(options)
	}
	return options
}

type Option func(*Options)

// WithOutputPath set the output path
func WithOutputPath(path string) Option {
	return func(options *Options) {
		options.OutputPath = path
	}
}

// WithAuthor set the author
func WithAuthor(author string) Option {
	return func(options *Options) {
		options.Author = author
	}
}

// WithProject set the project
func WithProject(project string) Option {
	return func(options *Options) {
		options.Project = project
	}
}

// WithOverride set the override flag
func WithOverride(override bool) Option {
	return func(options *Options) {
		options.Override = override
	}
}

// WithIgnoreExist set the ignore exist flag
func WithIgnoreExist(ignoreExist bool) Option {
	return func(options *Options) {
		options.IgnoreExist = ignoreExist
	}
}

// WithGoVersion set the go version
func WithGoVersion(version string) Option {
	return func(options *Options) {
		options.GoVersion = version
	}
}

// WithProtocPath set the protoc command path
func WithProtocPath(path string) Option {
	return func(options *Options) {
		options.ProtocPath = path
	}
}

// WithProtoFileIncludePath set the protoc command include .proto file paths
func WithProtoFileIncludePath(paths []string) Option {
	return func(options *Options) {
		options.ProtoFileIncludePath = paths
	}
}

// WithGoModule set the go module
func WithGoModule(module string) Option {
	return func(options *Options) {
		options.GoModule = module
	}
}

// WithAddr set the service address
func WithAddr(address string) Option {
	return func(options *Options) {
		options.Addr = address
	}
}

// WithPackageName set the package name
func WithPackageName(name string) Option {
	return func(options *Options) {
		options.PackageName = name
	}
}

// WithServiceName set the package name
func WithServiceName(name string) Option {
	return func(options *Options) {
		options.ServiceName = name
	}
}

// WithRPCName set the package name
func WithRPCName(name string) Option {
	return func(options *Options) {
		options.RPCName = name
	}
}

// WithRequestName set the package name
func WithRequestName(name string) Option {
	return func(options *Options) {
		options.RequestName = name
	}
}

// WithResponseName set the package name
func WithResponseName(name string) Option {
	return func(options *Options) {
		options.ResponseName = name
	}
}

// WithImportedPackageName set the imported package name
func WithImportedPackageName(name string) Option {
	return func(options *Options) {
		options.ImportedPackageName = name
	}
}

// WithServiceImplName set the service implementation name
func WithServiceImplName(name string) Option {
	return func(options *Options) {
		options.ServiceImplName = name
	}
}

// WithConfigPackagePath set the config package path
func WithConfigPackagePath(path string) Option {
	return func(options *Options) {
		options.ConfigPackagePath = path
	}
}

// WithNamespace set the namespace
func WithNamespace(namespace string) Option {
	return func(options *Options) {
		options.Namespace = namespace
	}
}

// WithVersion set the version
func WithVersion(version string) Option {
	return func(options *Options) {
		options.Version = version
	}
}

// WithPort set the port
func WithPort(port string) Option {
	return func(options *Options) {
		options.Port = port
	}
}

// WithTargetPort set the target port
func WithTargetPort(port string) Option {
	return func(options *Options) {
		options.TargetPort = port
	}
}

// WithImagePath set the image path
func WithImagePath(path string) Option {
	return func(options *Options) {
		options.ImagePath = path
	}
}
