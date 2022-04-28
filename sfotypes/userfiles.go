package sfotypes

type UserFile struct {
	FilePath string `json:"file_path"`
}

type UserData struct {
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
