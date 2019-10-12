package update

type state struct {
	Application       string `json:"app"`
	UpdateApplication string `json:"update"`
	Phase             Phase  `json:"phase"`
}
