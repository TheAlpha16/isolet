package event

type InfoOutput struct {
	Name       string `json:"name"`
	StartTime  int64  `json:"event_start"`
	EndTime    int64  `json:"event_end"`
	PostEvent  bool   `json:"post_event"`
	TeamLength int    `json:"team_length"`
}
