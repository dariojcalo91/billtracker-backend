package domain

type DashboardCounts struct {
	Total    int `json:"total"`
	Done     int `json:"done"`
	Upcoming int `json:"upcoming"`
	Overdue  int `json:"overdue"`
}

type BillDashboardStatus struct {
	Bill    *Bill    `json:"bill"`
	Payment *Payment `json:"payment,omitempty"`
}

type DashboardSummary struct {
	Month    string                 `json:"month"`
	Summary  DashboardCounts        `json:"summary"`
	Done     []*BillDashboardStatus `json:"done"`
	Upcoming []*BillDashboardStatus `json:"upcoming"`
	Overdue  []*BillDashboardStatus `json:"overdue"`
}
