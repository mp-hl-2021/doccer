package linter

type StubLinter struct {}

func (s * StubLinter) inspect(code string) (*InspectionResult, error) {
	return &InspectionResult{
		comments: "Text inspected",
	}, nil
}
