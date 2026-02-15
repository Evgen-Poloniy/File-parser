package entity

// Dataclass for writing information about parsed file in the database
type File struct {
	Filename string `db:"filename"`
}
