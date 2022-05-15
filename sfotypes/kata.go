package sfotypes

type KataResponse struct {
	ScansResult []Scans `json:"scans"`
}

type Scans struct {
	ScanId string `json:"scanId"`
	State  string `json:"state"`
}

type CheckQueueKata struct {
	Response     KataResponse `json:"response"`
	FileName     string       `json:"file_name"`
	IncomingFile FileToCheck  `json:"incoming_file"`
	ScanId       string       `json:"scan_id"`
}
