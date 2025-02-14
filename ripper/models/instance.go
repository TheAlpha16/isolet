package models

type Flag struct {
	FlagID     int64  `gorm:"primaryKey;autoIncrement;column:flagid" json:"-"`
	TeamID     int64  `gorm:"not null;column:teamid" json:"-"`
	ChallID    int    `gorm:"not null;column:chall_id" json:"chall_id"`
	Flag       string `gorm:"not null" json:"-"`
	Password   string `gorm:"type:text" json:"password"`
	Port       int    `gorm:"type:integer" json:"port"`
	Hostname   string `gorm:"type:text" json:"hostname"`
	Deadline   int64  `gorm:"type:integer" json:"deadline"`
	Extended   int    `gorm:"type:integer;default:1" json:"-"`
	Deployment string `gorm:"-" json:"deployment"`
}

type Running struct {
	RunID   int64 `gorm:"primaryKey;autoIncrement;column:runid" json:"-"`
	TeamID  int64 `gorm:"not null;column:teamid" json:"teamid"`
	ChallID int   `gorm:"not null;column:chall_id" json:"chall_id"`
}
