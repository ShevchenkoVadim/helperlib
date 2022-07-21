package sfotypes

type Configurations struct {
	SrcPath            string       `yaml:"srcpath"`
	DestPath           string       `yaml:"destpath"`
	Separator          string       `yaml:"separator"`
	LogFilePath        string       `yaml:"logfilepath"`
	LogQueue           string       `yaml:"logqueue"`
	LogDateFormat      string       `yaml:"logDateFormat"`
	DirFiles           string       `yaml:"dirfiles"`
	MaxSizeLog         int          `yaml:"maxsizelog"`
	SSLCert            Certificate  `yaml:"ssl_cert"`
	MQ                 MqConfig     `yaml:"mq"`
	MS                 MSConfig     `yaml:"ms"`
	IW                 IWConfig     `yaml:"iw"`
	Kata               KataConfig   `yaml:"kata"`
	Notify             NotifyConfig `yaml:"notify"`
	DBConn             DBConnection `yaml:"db_conn"`
	DebugEnable        bool         `yaml:"debugenable"`
	Hostname           string       `yaml:"hostname"`
	Debug              bool         `yaml:"debug"`
	Filter             string       `yaml:"filter"`
	FilterExcludeFiles string       `yaml:"filter_exclude_files"`
	ScanPath           string       `yaml:"scanpath"`
	OsType             string       `yaml:"ostype"`
	ScanTimeout        int          `yaml:"scantimeout"`
	MoverType          int          `yaml:"movertype"`
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
	Notifier        string `yaml:"notifier"`
	//LogQueue    string `yaml:"logqueue"`
}

type MSConfig struct {
	ApiUrl        string `yaml:"api_url"`
	Timeout       int    `yaml:"timeout"`
	AnalysisDepth int    `yaml:"analysis_depth"`
	XApiKey       string `yaml:"x_api_key"`
}

type DBConnection struct {
	Server   string `yaml:"server"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Db       string `yaml:"db"`
	Timeout  int    `yaml:"timeout"`
}

type KataConfig struct {
	ApiUri   string `yaml:"api"`
	KataId   string `yaml:"kataid"`
	Instance string `yaml:"instance"`
	Timeout  int    `yaml:"timeout"`
}

type IWConfig struct {
	ApiUrl        string `yaml:"api_url"`
	Token         string `yaml:"token"`
	SystemId      string `yaml:"system_id"`
	SystemClass   string `yaml:"system_class"`
	SystemService string `yaml:"system_service"`
}

type NotifyConfig struct {
	SMPTHost           string `yaml:"smpt_host"`
	SMPTPort           string `yaml:"smpt_port"`
	MailDomain         string `yaml:"mail_domain"`
	FromUser           string `yaml:"from_user"`
	SendMail           bool   `yaml:"send_mail"`
	CitrixBrocker      string `yaml:"citrix_srv"`
	CitrixDomain       string `yaml:"citrix_domain"`
	CitrixMessageTitle string `yaml:"citrix_message_title"`
}
