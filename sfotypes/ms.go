package sfotypes

type MsResponce struct {
	Errors []ErrorsMessage `json:"errors"`
	Data   DataResponce    `json:"data"`
	ScanId string          `json:"scan_id"`
}

type DataResponce struct {
	FileUri string `json:"file_uri"`
	Ttl     int    `json:"ttl"`
	Results Result `json:"result"`
}

type Task struct {
	FileUri     string      `json:"file_uri"`
	FileName    string      `json:"file_name"`
	ShortResult bool        `json:"short_result"`
	Options     TaskOptions `json:"options"`
}

type TaskOptions struct {
	AnalysisDepth      int      `json:"analysis_depth"`
	PasswordsForUnpack []string `json:"passwords_for_unpack"`
}

type CheckQueue struct {
	Response     MsResponce  `json:"response"`
	FileName     string      `json:"file_name"`
	IncomingFile FileToCheck `json:"incoming_file"`
}

type Result struct {
	ScanState string          `json:"scan_state"`
	Duration  float64         `json:"duration"`
	Verdict   string          `json:"verdict"`
	Errors    []ErrorsMessage `json:"errors"`
}

type ErrorsMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
