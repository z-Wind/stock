package main

import (
	"bytes"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
)

func getCurExePath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", errors.Wrapf(err, "exec.LookPath")
	}

	//得到全路径，比如在windows下E:\\golang\\test\\a.exe
	path, err := filepath.Abs(file)
	if err != nil {
		return "", errors.Wrapf(err, "filepath.Abs")
	}

	rst := filepath.Dir(path)

	return rst, nil
}

func getCurScriptPath() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("runtime.Caller Fail")
	}

	//得到全路径，比如在windows下E:\\golang\\test\\a.exe
	path, err := filepath.Abs(file)
	if err != nil {
		return "", errors.Wrapf(err, "filepath.Abs")
	}

	rst := filepath.Dir(path)

	return rst, nil
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func removeDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]struct{}{}
	result := []string{}

	for v := range elements {
		if _, ok := encountered[elements[v]]; ok {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = struct{}{}
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func parseTemplate(fileName string, data interface{}) (output []byte, err error) {
	var buf bytes.Buffer
	template, err := template.ParseFiles(fileName)
	if err != nil {
		return nil, err
	}
	err = template.Execute(&buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
