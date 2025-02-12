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

type Instance struct {
	ChallID    int    `gorm:"column:chall_id" json:"chall_id"`
	Password   string `gorm:"column:password" json:"password"`
	Port       int    `gorm:"column:port" json:"-"`
	Hostname   string `gorm:"column:hostname" json:"-"`
	Deadline   int64  `gorm:"column:deadline" json:"deadline"`
	Deployment string `gorm:"column:deployment" json:"-"`
	ConnString string `gorm:"-" json:"connstring"`
}

type Image struct {
	IID        int    `gorm:"primaryKey;autoIncrement;column:iid" json:"iid"`
	ChallID    int    `gorm:"not null;column:chall_id" json:"chall_id"`
	Image      string `gorm:"not null" json:"image"`
	Deployment string `gorm:"type:deployment_type;default:http" json:"deployment"`
	Port       int    `gorm:"default:80" json:"port"`
	Subd       string `gorm:"default:localhost" json:"subd"`
	CPU        int    `gorm:"default:5;column:cpu" json:"-"`
	Memory     int    `gorm:"default:10;column:mem" json:"-"`
}

type ExtendDeadline struct {
	Deadline int64 `json:"deadline"`
}

func (Running) TableName() string {
	return "running"
}
