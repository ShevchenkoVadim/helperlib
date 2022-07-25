package sfotypes

type LogMsg struct {
	HostName    string `json:"hostName"`
	ServiceName string `json:"serviceName"`
	Msg         string `json:"msg"`
	TimeStamp   int64  `json:"timeStamp"`
}
