package helm

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type KubeconformResult struct {
	Chart  string
	Values string
	Output string
	Err    error
}

func (helm *Helm) Kubeconform(chartPath string, helmOptions *TemplateOptions, kubeconformOptions *KubeconformOptions) (kubeconformResult []*KubeconformResult, err error) {
	_, err = helm.Dependency(chartPath)
	if err != nil {
		return nil, err
	}

	valuesFiles, err := helm.ListValuesFiles(chartPath)
	if err != nil {
		return nil, err
	}

	resultsChan := make(chan *KubeconformResult, len(valuesFiles))
	var wg sync.WaitGroup
	for _, valuesFile := range valuesFiles {
		wg.Add(1)
		go func(chartPath string, valuesFile string) {
			defer wg.Done()
			helmOptions.ValuesFile = []string{valuesFile}
			output, err := helm.Template(chartPath, helmOptions.GetArgs()...)
			if err != nil {
				resultsChan <- &KubeconformResult{
					Chart:  filepath.Base(chartPath),
					Values: filepath.Base(valuesFile),
					Output: output,
					Err:    err,
				}
				return
			}

			output, err = kubeconform(output, kubeconformOptions.GetArgs()...)
			resultsChan <- &KubeconformResult{
				Chart:  filepath.Base(chartPath),
				Values: filepath.Base(valuesFile),
				Output: output,
				Err:    err,
			}
		}(chartPath, valuesFile)
	}

	go func() {
		defer close(resultsChan)
		wg.Wait()
	}()

	for result := range resultsChan {
		kubeconformResult = append(kubeconformResult, result)
	}

	return kubeconformResult, nil
}

func kubeconform(input string, options ...string) (output string, err error) {
	var stdout bytes.Buffer

	subProcess := exec.Command("kubeconform", options...)
	subProcess.Stdout = &stdout
	subProcess.Stderr = &stdout
	subProcess.Stdin = strings.NewReader(input)

	err = subProcess.Run()
	if err != nil {
		return strings.Trim(stdout.String(), "\n"), err
	}

	return strings.Trim(stdout.String(), "\n"), nil
}
