package helmunit

import (
	"encoding/json"
	"fmt"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	v "helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/releaseutil"
	"regexp"
	"sigs.k8s.io/yaml"
)

func splitChart(chart string) map[string][]byte {
	splitChart := releaseutil.SplitManifests(chart)

	nameRegex := regexp.MustCompile("# Source: [^/]+/(.+)")
	files := make(map[string][]byte)

	for _, file := range splitChart {
		submatch := nameRegex.FindStringSubmatch(file)
		if len(submatch) == 0 {
			continue
		}

		name := submatch[1]

		files[name] = []byte(file)
	}

	return files
}

func Template(name, namespace, chartPath, filePath string, valueFilePaths, overrideValues []string, output interface{}) error  {
	client := defaultClient(name, namespace)

	p := getter.All(&cli.EnvSettings{})
	valueOpts := &v.Options{
		ValueFiles:   valueFilePaths,
		Values:       overrideValues,
	}

	values, err := valueOpts.MergeValues(p)
	if err != nil {
		return err
	}
	chart, err := loader.Load(chartPath)
	if err != nil {
		return err
	}
	release, err := client.Run(chart, values)
	if err != nil {
		return err
	}

	manifests := splitChart(release.Manifest)

	if _, exists := manifests[filePath]; !exists {
		return fmt.Errorf("no file found at path: %s", filePath)
	}

	jsonBytes, err := yaml.YAMLToJSON(manifests[filePath])
	if err != nil {
		return err
	}

	if err = json.Unmarshal(jsonBytes, &output); err != nil {
		return err
	}

	return nil
}

func defaultClient(name, namespace string) *action.Install {
	client := action.NewInstall(&action.Configuration{})
	client.Version = ">0.0.0-0"
	client.ReleaseName = name
	client.Namespace = namespace
	client.ClientOnly = true
	client.DryRun = true

	return client
}

