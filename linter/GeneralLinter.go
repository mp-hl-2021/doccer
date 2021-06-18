package linter

import (
	"doccer/data"
)

type GeneralLinter struct {
	mapper map[string]Linter
}

func NewGeneralLinter() GeneralLinter {
	return GeneralLinter {make(map[string]Linter) }
}

func (g *GeneralLinter) RegisterNewLinter(langName string, linter Linter) {
	g.mapper[langName] = linter
}

func (g *GeneralLinter) CheckCode(doc data.Doc) data.Doc {
	linter, ok := g.mapper[doc.Lang]
	if !ok {
		doc.LinterStatus = "No inspection for " + doc.Lang
		return doc
	}
	lintRes, err := linter.inspect(doc.Text)
	if err != nil {
		doc.LinterStatus = "No inspection"
		return doc
	}
	doc.LinterStatus = lintRes.comments
	return doc
}
