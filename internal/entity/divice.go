package entity

// Dataclass for writing data from parsed files in the database
type Device struct {
	UnitGUID string `db:"unit_guid"`
	InvID    string `db:"inv_id"`
}
