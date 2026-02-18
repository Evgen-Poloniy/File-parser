package dto

// Response represents a data record from database
// @Description Data record structure returned from the database
type Response struct {
	N         *int    `json:"n"`
	UnitGUID  string  `json:"unit_guid"`
	InvID     string  `json:"inv_id"`
	MQTT      *string `json:"mqtt"`
	MsgID     *string `json:"msg_id"`
	Text      *string `json:"text"`
	Context   *string `json:"context"`
	Class     *string `json:"class"`
	Level     *int    `json:"level"`
	Area      *string `json:"area"`
	Addr      *string `json:"addr"`
	Block     *string `json:"block"`
	Type      *string `json:"type"`
	Bit       *int    `json:"bit"`
	InvertBit *int    `json:"invert_bit"`
}

// FileErrorInfo represents parsing error information
// @Description Information about errors that occurred during file parsing
type FileErrorInfo struct {
	Filename string  `json:"filename"`
	Line     *int    `json:"line"`
	LineData *string `json:"line_data"`
	ErrorMsg string  `json:"error_msg"`
}
