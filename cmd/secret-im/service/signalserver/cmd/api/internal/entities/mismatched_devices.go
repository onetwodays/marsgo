package entities
// 不匹配设备
type MismatchedDevices struct {
	MissingDevices []int64 `json:"missingDevices"`
	ExtraDevices   []int64 `json:"extraDevices"`
}

