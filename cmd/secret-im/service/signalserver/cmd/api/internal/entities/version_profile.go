package entities

type  VersionedProfile struct {
	Version string `json:"version"`
	Name string `json:"name"`
	Avatar string `json:"avatar"`
	Commitment string `json:"commitment"`
}
