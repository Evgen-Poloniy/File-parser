package entity

// Dataclass for writing data from parsed files in the database
type DeviceData struct {
	N         *int    `db:"n"`
	MQTT      *string `db:"mqtt"`
	MsgID     string  `db:"msg_id"`
	Text      string  `db:"text"`
	Context   string  `db:"context"`
	Class     string  `db:"class"`
	Level     *int    `db:"level"`
	Area      string  `db:"area"`
	Addr      string  `db:"addr"`
	Block     *string `db:"block"`
	Type      *string `db:"type"`
	Bit       *int    `db:"bit"`
	InvertBit *int    `db:"invert_bit"`
	UnitGUID  string  `db:"unit_guid"`
}
