package sfotypes

type LogMsg struct {
	ServiceName string `json:"serviceName"`
	Msg         string `json:"msg"`
	TimeStamp   int64  `json:"timeStamp"`
}
