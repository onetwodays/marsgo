package criticality

// Criticality is
type Criticality string

// criticality:关键性?
var (
	// EmptyCriticality is used to mark any invalid criticality,
	//and the empty criticality will be parsed as the default criticality later.
	EmptyCriticality = Criticality("")
	// CriticalPlus is reserved for the most critical requests,
	//those that will result in serious user-visible impact if they fail.
	CriticalPlus = Criticality("CRITICAL_PLUS")
	// Critical is the default value for requests sent from production jobs.
	//These requests will result in user-visible impact, but the impact may be less severe than those of CRITICAL_PLUS.
	//Services are expected to provision enough capacity for all expected CRITICAL and CRITICAL_PLUS traffic.
	//从生产作业发送的请求的默认值。
	//这些请求将导致用户可见的影响，但影响可能不如CRITICAL_PLUS严重。 服务应为所有预期的CRITICAL和CRITICAL_PLUS流量提供足够的容量。
	Critical = Criticality("CRITICAL")
	// SheddablePlus is traffic for which partial unavailability is expected. This is the default for batch jobs, which can retry requests minutes or even hours later.
	SheddablePlus = Criticality("SHEDDABLE_PLUS") //SheddablePlus是预期部分不可用的流量,这是批处理作业的默认设置，可以在几分钟甚至几小时后重试请求。
	// Sheddable is traffic for which frequent partial unavailability and occasional full unavailability is expected.
	Sheddable = Criticality("SHEDDABLE") //可预期流量是指预计经常出现部分不可用和偶尔完全不可用的流量

	// higher is more critical
	_criticalityEnum = map[Criticality]int{
		CriticalPlus:  40,
		Critical:      30,
		SheddablePlus: 20,
		Sheddable:     10,
	}

	_defaultCriticality = Critical
)

// Value is used to get criticality value, higher value is more critical.
func Value(in Criticality) int {
	v, ok := _criticalityEnum[in]
	if !ok {
		return _criticalityEnum[_defaultCriticality]
	}
	return v
}

// Higher will compare the input criticality with self, return true if the input is more critical than self.
func (c Criticality) Higher(in Criticality) bool {
	return Value(in) > Value(c)
}

// Parse will parse raw criticality string as valid critality. Any invalid input will parse as empty criticality.
func Parse(raw string) Criticality {
	crtl := Criticality(raw)
	if _, ok := _criticalityEnum[crtl]; ok {
		return crtl
	}
	return EmptyCriticality
}

// Exist is used to check criticality is exist in several enumeration.
func Exist(c Criticality) bool {
	_, ok := _criticalityEnum[c]
	return ok
}
