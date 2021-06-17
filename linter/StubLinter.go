package linter

type StubLinter struct {}

func (s * StubLinter) inspect(code string) (*InspectionResult, error) {
	return &InspectionResult{
		codeCorrect: true,
		comments: "Text inspected",
	}, nil
}
