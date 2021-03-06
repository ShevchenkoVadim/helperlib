package sfotypes

type UserFile struct {
	FilePath string `json:"file_path"`
}

type UserData struct {
	Service   string     `json:"service"`
	UserName  string     `json:"user_name"`
	Files     []UserFile `json:"files"`
	TimeStamp int64      `json:"timeStamp"`
}

type FileToCheck struct {
	UserName  string `json:"username"`
	UserFile  string `json:"userfile"`
	HostName  string `json:"hostname"`
	Timestamp int64  `json:"timestamp"`
}

type ServiceTask struct {
	Service   string      `json:"service"`
	Hostname  string      `json:"hostname"`
	File      FileToCheck `json:"file"`
	UserDatas UserData    `json:"user_data"`
	Status    string      `json:"status"`
	StatusExt string      `json:"status_ext"`
}
