package sfotypes

type Configurations struct {
	SrcPath     string      `yaml:"srcpath"`
	DestPath    string      `yaml:"destpath"`
	Separator   string      `yaml:"separator"`
	LogFilePath string      `yaml:"logfilepath"`
	DirFiles    string      `yaml:"dirfiles"`
	MaxSizeLog  int         `yaml:"maxsizelog"`
	SSLCert     Certificate `yaml:"ssl_cert"`
	MQ          MqConfig    `yaml:"mq"`
	MS          MSConfig    `yaml:"ms"`
	DBConn      string      `yaml:"db_conn"`
	DebugEnable bool        `yaml:"debugenable"`
	Hostname    string      `yaml:"hostname""`
}

type Certificate struct {
	SslCA  string `yaml:"ssl_ca"`
	SslPem string `yaml:"ssl_pem"`
	SslKey string `yaml:"ssl_key"`
}

type MqConfig struct {
	Url         string `yaml:"url"`
	SrcQueue    string `yaml:"srcqueue"`
	ManageQueue string `yaml:"managequeue"`
	LogQueue    string `yaml:"logqueue"`
}

type MSConfig struct {
	ApiUrl        string `yaml:"api_url"`
	Timeout       int    `yaml:"timeout"`
	AnalysisDepth int    `yaml:"analysis_depth"`
	XApiKey       string `yaml:"x_api_key"`
}
