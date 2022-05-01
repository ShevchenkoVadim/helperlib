package sfotypes

type Configurations struct {
	SrcPath       string      `yaml:"srcpath"`
	DestPath      string      `yaml:"destpath"`
	Separator     string      `yaml:"separator"`
	LogFilePath   string      `yaml:"logfilepath"`
	LogQueue      string      `yaml:"logqueue"`
	LogDateFormat string      `yaml:"logDateFormat"`
	DirFiles      string      `yaml:"dirfiles"`
	MaxSizeLog    int         `yaml:"maxsizelog"`
	SSLCert       Certificate `yaml:"ssl_cert"`
	MQ            MqConfig    `yaml:"mq"`
	MS            MSConfig    `yaml:"ms"`
	DBConn        string      `yaml:"db_conn"`
	DebugEnable   bool        `yaml:"debugenable"`
	Hostname      string      `yaml:"hostname"`
	Debug         bool        `yaml:"debug"`
	Filter        string      `yaml:"filter"`
	ScanPath      string      `yaml:"scanpath"`
	OsType        string      `yaml:"ostype"`
	ScanTimeout   int         `yaml:"scantimeout"`
}

type Certificate struct {
	SslCA  string `yaml:"ssl_ca"`
	SslPem string `yaml:"ssl_pem"`
	SslKey string `yaml:"ssl_key"`
}

type MqConfig struct {
	Url             string `yaml:"url"`
	SrcQueue        string `yaml:"srcqueue"`
	ManageQueue     string `yaml:"managequeue"`
	TempQueue       string `yaml:"tempqueue"`
	MSQueue         string `yaml:"msqueue"`
	KataQueue       string `yaml:"kataqueue"`
	InfoWatchQueue  string `yaml:"infowatchqueue"`
	QuarantineQueue string `yaml:"quarantinequeue"`
	DestQueue       string `yaml:"destqueue"`
	//LogQueue    string `yaml:"logqueue"`
}

type MSConfig struct {
	ApiUrl        string `yaml:"api_url"`
	Timeout       int    `yaml:"timeout"`
	AnalysisDepth int    `yaml:"analysis_depth"`
	XApiKey       string `yaml:"x_api_key"`
}
