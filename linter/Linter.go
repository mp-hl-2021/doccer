package linter

type Linter interface {
	inspect(code string) (*InspectionResult, error)
}

type InspectionResult struct {
	codeCorrect bool
	comments string
}
