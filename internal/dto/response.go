package dto

// Dataclass for response from database
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

// Dataclass for response with information about parse error from database
type FileErrorInfo struct {
	Filename string  `json:"filename"`
	Line     *int    `json:"line"`
	LineData *string `json:"line_data"`
	ErrorMsg string  `json:"error_msg"`
}
