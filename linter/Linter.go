package linter

type Linter interface {
	inspect(code string) (*InspectionResult, error)
}

type InspectionResult struct {
	comments string
}
