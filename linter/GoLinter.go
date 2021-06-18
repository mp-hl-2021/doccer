package linter

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type GoLinter struct {}

func (s * GoLinter) inspect(code string) (*InspectionResult, error) {
	file, err := ioutil.TempFile("", "tmp*.go")
	if err != nil {
		return nil, err
	}
	defer os.Remove(file.Name())
	_, err = file.WriteString(code)
	if err != nil {
		return nil, err
	}

	cmd, err := exec.Command("staticcheck", file.Name()).CombinedOutput()

	res := string(cmd)
	res = strings.Replace(res, file.Name(), "", -1)
	res = strings.TrimSpace(res)
	if len(res) == 0 {
		res = "OK"
	}
	return &InspectionResult{
		comments: res,
	}, nil
}