package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/philippe-vandermoere/helm-plugin-kubeconform/pkg/helm"
)

var version = "undefined"

func main() {
	chartPath, helmTemplateOptions, kubeconformOptions, err := parseArgs()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	output, err := helm.New().Kubeconform(*chartPath, helmTemplateOptions, kubeconformOptions)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	exitCode := 0
	for _, result := range output {
		fmt.Printf("validate chart \"%s\" with values \"%s\":\n%s\n", result.Chart, result.Values, result.Output)
		if result.Err != nil {
			exitCode = 1
		}
	}

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(exitCode)
	}
}

func parseArgs() (chart *string, helmTemplateOptions *helm.TemplateOptions, kubeconformOptions *helm.KubeconformOptions, err error) {
	if len(os.Args) < 2 {
		return nil, nil, nil, errors.New("\"helm-kubeconform\" requires at least 1 argument")
	}

	helmTemplateOptions = &helm.TemplateOptions{}
	kubeconformOptions = &helm.KubeconformOptions{}
	chart = &os.Args[1]
	helmKubeconform := flag.NewFlagSet("helm-kubeconform", flag.ContinueOnError)

	// options
	help := helmKubeconform.BoolP("help", "h", false, "show this help message and exit")
	argVersion := helmKubeconform.BoolP("version", "v", false, "print the version information")

	// helm template options
	helmKubeconform.StringArrayVarP(&helmTemplateOptions.ValuesFile, "values", "f", []string{}, "specify values in a YAML file or a URL (can specify multiple)")
	helmKubeconform.StringVar(&helmTemplateOptions.KubernetesVersion, "kubernetes-version", "master", "version of Kubernetes to validate against, e.g.: 1.18.0")
	helmKubeconform.StringVarP(&helmTemplateOptions.Namespace, "namespace", "n", "", "namespace scope for this request")

	// kubeconform options
	helmKubeconform.IntVar(&kubeconformOptions.Goroutines, "goroutines", 4, "number of goroutines to run concurrently")
	helmKubeconform.StringVar(&kubeconformOptions.Output, "output", "text", "output format - json, junit, tap, text")
	helmKubeconform.StringVar(&kubeconformOptions.Reject, "reject", "", "comma-separated list of kinds or GVKs to reject)")
	helmKubeconform.StringArrayVar(&kubeconformOptions.SchemaLocation, "schema-location", []string{"default"}, "override schemas location search path (can be specified multiple times)")
	helmKubeconform.StringVar(&kubeconformOptions.Skip, "skip", "", "comma-separated list of kinds or GVKs to ignore")
	helmKubeconform.BoolVar(&kubeconformOptions.Strict, "strict", true, "disallow additional properties not in schema or duplicated keys")
	helmKubeconform.BoolVar(&kubeconformOptions.Summary, "summary", true, "print a summary at the end (ignored for junit output)")
	helmKubeconform.BoolVar(&kubeconformOptions.Verbose, "verbose", false, "print results for all resources (ignored for tap and junit output)")

	err = helmKubeconform.Parse(os.Args[1:])
	if err != nil {
		return nil, nil, nil, err
	}

	if *help {
		usage(helmKubeconform)
		os.Exit(0)
	}

	if *argVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	// copy kubernetes version form helm option to kubeconform options
	kubeconformOptions.KubernetesVersion = helmTemplateOptions.KubernetesVersion

	return chart, helmTemplateOptions, kubeconformOptions, nil
}

func usage(flagSet *flag.FlagSet) {
	options := []*flag.Flag{}
	helmTemplate := []*flag.Flag{}
	kubeconform := []*flag.Flag{}

	flagSet.VisitAll(func(f *flag.Flag) {
		switch f.Name {
		case "values", "kubernetes-version", "namespace":
			helmTemplate = append(helmTemplate, f)
		case "goroutines", "output", "reject", "schema-location", "skip", "strict", "summary", "verbose":
			kubeconform = append(kubeconform, f)
		default:
			options = append(options, f)
		}
	})

	lines := []string{
		"helm kubeconform chart [flags]",
		"\npositional arguments:",
		"  chart\x00chart",
	}
	lines = append(lines, usageSection("options:", options)...)
	lines = append(lines, usageSection("Helm template options:", helmTemplate)...)
	lines = append(lines, usageSection("Kubeconform options:", kubeconform)...)

	maxlen := 0
	for _, line := range lines {
		if strings.Contains(line, "\x00") {
			split := strings.Split(line, "\x00")
			if len(split[0]) > maxlen {
				maxlen = len(split[0])
			}
		}
	}

	for _, line := range lines {
		sidx := strings.Index(line, "\x00")
		if sidx > 0 {
			spacing := strings.Repeat(" ", maxlen-sidx)
			fmt.Println(line[:sidx], spacing, strings.ReplaceAll(line[sidx+1:], "\n", "\n"+strings.Repeat(" ", maxlen+2)))
		} else {
			fmt.Println(line)
		}
	}
}

func usageSection(title string, flags []*flag.Flag) (lines []string) {
	for i, flag := range flags {
		if i == 0 {
			lines = append(lines, fmt.Sprintf("\n%s", title))
		}
		lines = append(lines, formatFlag(flag))
	}
	return
}

func formatFlag(f *flag.Flag) (line string) {
	if f.Hidden {
		return line
	}
	if f.Shorthand != "" && f.ShorthandDeprecated == "" {
		line = fmt.Sprintf("  -%s, --%s", f.Shorthand, f.Name)
	} else {
		line = fmt.Sprintf("      --%s", f.Name)
	}
	varname, usage := flag.UnquoteUsage(f)
	if varname != "" {
		line += " " + varname
	}
	if f.NoOptDefVal != "" {
		switch f.Value.Type() {
		case "string":
			line += fmt.Sprintf("[=\"%s\"]", f.NoOptDefVal)
		case "bool":
			if f.NoOptDefVal != "true" {
				line += fmt.Sprintf("[=%s]", f.NoOptDefVal)
			}
		case "count":
			if f.NoOptDefVal != "+1" {
				line += fmt.Sprintf("[=%s]", f.NoOptDefVal)
			}
		default:
			line += fmt.Sprintf("[=%s]", f.NoOptDefVal)
		}
	}
	line += "\x00" + usage
	if f.Value.Type() == "string" {
		if f.DefValue != "" {
			line += fmt.Sprintf(" (default %q)", f.DefValue)
		}
	} else if f.Value.Type() != "bool" {
		line += fmt.Sprintf(" (default %s)", f.DefValue)
	}

	if len(f.Deprecated) != 0 {
		line += fmt.Sprintf(" (DEPRECATED: %s)", f.Deprecated)
	}

	return line
}
