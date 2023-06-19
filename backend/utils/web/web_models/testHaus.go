package webmodels

type QueryParameter struct {
	StartDate string
	EndDate   string
	Board     string
	Reason    string
	Name      string
	Result    string
	OrderBy   string
}

type QueryValue struct {
	Filter    string
	Value     string
	Condition string
}
