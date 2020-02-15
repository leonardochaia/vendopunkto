package conf

type Startup struct {
	Logger  StartupLogger  `mapstructure:"logger"`
	PIDFile string         `mapstructure:"pidfile"`
	Server  StartupServer  `mapstructure:"server"`
	Storage StartupStorage `mapstructure:"storage"`
}

type StartupLogger struct {
	Level string `mapstructure:"level"`
}

type StartupServer struct {
	LogRequests     bool                  `mapstructure:"log_requests"`
	ProfilerEnabled bool                  `mapstructure:"profiler_enabled"`
	ProfilerPath    string                `mapstructure:"profiler_path"`
	Internal        StartupInternalServer `mapstructure:"internal"`
	Public          StartupPublicServer   `mapstructure:"public"`
}

type StartupInternalServer struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type StartupPublicServer struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type StartupStorage struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"sslmode"`
	Type     string `mapstructure:"type"`
}
