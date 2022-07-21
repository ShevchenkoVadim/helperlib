package sfotypes

type TextData struct {
	FileScanerMessage  string `yaml:"file_scaner_message"`
	QuarantineMessage  string `yaml:"quarantine_message"`
	DestinationMessage string `yaml:"destination_message"`
}
