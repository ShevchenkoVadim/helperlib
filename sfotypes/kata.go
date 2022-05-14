package sfotypes

type KataResponse struct {
	ScansResult []Scans `json:"scans"`
}

type Scans struct {
	ScanId string `json:"scanId"`
	State  string `json:"state"`
}
