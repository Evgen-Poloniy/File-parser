package entity

// Dataclass for writing data from parsed files in the database
// Unfortunately, tegs is not used because of using standard library sql
type Device struct {
	InvID    string `db:"inv_id"`
	UnitGUID string `db:"unit_guid"`
}
