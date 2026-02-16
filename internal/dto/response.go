package dto

// Dataclass for response for client
type Response struct {
	N         *int    `json:"n"`
	UnitGUID  string  `json:"unit_guid"`
	InvID     string  `json:"inv_id"`
	MQTT      *string `json:"mqtt"`
	MsgID     string  `json:"msg_id"`
	Text      string  `json:"text"`
	Context   string  `json:"context"`
	Class     string  `json:"class"`
	Level     *int    `json:"level"`
	Area      string  `json:"area"`
	Addr      string  `json:"addr"`
	Block     *string `json:"block"`
	Type      *string `json:"type"`
	Bit       *int    `json:"bit"`
	InvertBit *int    `json:"invert_bit"`
}
