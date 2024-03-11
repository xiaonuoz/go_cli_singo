package report_sub

// Order 订单
type ReportSub struct {
	Id              uint64
	Name            string
	ReportType      uint
	ReportForm      uint
	CronType        uint
	CronTime        string
	Email           bool
	SendOffiaccount bool
	CreatedAt       int64
	UpdatedAt       int64
	MonitorIds      string
	WafIds          string
	NodsIds         string // 云墙ids
	Uuid            string
	Pid             string
}
