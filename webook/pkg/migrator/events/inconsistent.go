package events

type InconsistentEvent struct {
	ID        int64
	Direction string
	Type      string
}

const (
	// InconsistentEventTypeTargetMissing 校验的目标数据，缺了这一条
	InconsistentEventTypeTargetMissing = "target_missing"
	// InconsistentEventTypeNEQ 不相等
	InconsistentEventTypeNEQ = "neq"
	// InconsistentEventTypeBaseMissing 校验的源数据，缺了这一条（反向校验用）
	InconsistentEventTypeBaseMissing = "base_missing"
)
