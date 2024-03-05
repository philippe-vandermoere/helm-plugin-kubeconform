package helm

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
)

type Helm struct{}

func New() *Helm {
	return &Helm{}
}

func (helm *Helm) Dependency(chartPath string, options ...string) (output string, err error) {
	cmd := append([]string{"dependency", "update", chartPath}, options...)
	output, err = helm.execute(cmd...)
	if err != nil {
		return "", err
	}

	return output, err
}

func (helm *Helm) Template(chartPath string, options ...string) (output string, err error) {
	cmd := append([]string{"template", chartPath}, options...)
	output, err = helm.execute(cmd...)
	if err != nil {
		return "", err
	}

	return output, err
}

func (helm *Helm) ListValuesFiles(chartPath string) (valuesFiles []string, err error) {
	valuesFiles, err = filepath.Glob(chartPath + "/" + helm.getValuesPath() + "/" + helm.getValuesFilePattern())
	if err != nil {
		return []string{}, err
	}

	return valuesFiles, nil
}

func (helm *Helm) execute(args ...string) (output string, err error) {
	var stdout bytes.Buffer

	cmd := getEnv("HELM_BIN", "helm")
	subProcess := exec.Command(cmd, args...)
	subProcess.Stdout = &stdout
	subProcess.Stderr = &stdout

	err = subProcess.Run()
	if err != nil {
		return stdout.String(), err
	}

	return stdout.String(), nil
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func (helm *Helm) getValuesPath() string {
	return "ci"
}

func (helm *Helm) getValuesFilePattern() string {
	return "*-values.yaml"
}
