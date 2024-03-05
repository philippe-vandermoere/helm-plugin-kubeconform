package helm

import (
	"strconv"
)

type TemplateOptions struct {
	ValuesFile        []string
	KubernetesVersion string
	Namespace         string
}

func (helmTemplateOptions *TemplateOptions) GetArgs() (options []string) {
	if helmTemplateOptions.KubernetesVersion != "master" {
		options = append(options, "--kube-version", helmTemplateOptions.KubernetesVersion)
	}
	for _, value := range helmTemplateOptions.ValuesFile {
		options = append(options, "--values", value)
	}
	if helmTemplateOptions.Namespace != "" {
		options = append(options, "--namespace", helmTemplateOptions.Namespace)
	}

	return options
}

type KubeconformOptions struct {
	Cache               string
	ExitOnError         bool
	IgnoreMissingSchema bool
	SkipTLSVerify       bool
	KubernetesVersion   string
	Goroutines          int
	Output              string
	Reject              string
	SchemaLocation      []string
	Skip                string
	Strict              bool
	Summary             bool
	Verbose             bool
}

func (kubeconformOptions *KubeconformOptions) GetArgs() (options []string) {
	if kubeconformOptions.Cache != "" {
		options = append(options, "-cache", kubeconformOptions.Cache)
	}
	if kubeconformOptions.ExitOnError {
		options = append(options, "-exit-on-error")
	}
	if kubeconformOptions.IgnoreMissingSchema {
		options = append(options, "-ignore-missing-schemas")
	}
	if kubeconformOptions.SkipTLSVerify {
		options = append(options, "-skip-tls-verify")
	}
	if kubeconformOptions.KubernetesVersion != "master" {
		options = append(options, "-kubernetes-version", kubeconformOptions.KubernetesVersion)
	}
	if kubeconformOptions.Goroutines != 4 {
		options = append(options, "-n", strconv.Itoa(kubeconformOptions.Goroutines))
	}
	if kubeconformOptions.Output != "text" {
		options = append(options, "-output", kubeconformOptions.Output)
	}
	if kubeconformOptions.Reject != "" {
		options = append(options, "-reject", kubeconformOptions.Reject)
	}
	for _, value := range kubeconformOptions.SchemaLocation {
		options = append(options, "-schema-location", value)
	}
	if kubeconformOptions.Skip != "" {
		options = append(options, "-skip", kubeconformOptions.Skip)
	}
	if kubeconformOptions.Strict {
		options = append(options, "-strict")
	}
	if kubeconformOptions.Summary {
		options = append(options, "-summary")
	}
	if kubeconformOptions.Verbose {
		options = append(options, "-verbose")
	}

	return options
}
