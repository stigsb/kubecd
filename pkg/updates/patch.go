package updates

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

func GroupImageUpdatesByReleasesFile(imageUpdates []ImageUpdate) map[string][]ImageUpdate {
	updatesPerFile := make(map[string][]ImageUpdate)
	for _, update := range imageUpdates {
		file := update.Release.FromFile
		if _, found := updatesPerFile[file]; !found {
			updatesPerFile[file] = make([]ImageUpdate, 0)
		}
		updatesPerFile[file] = append(updatesPerFile[file], update)
	}
	return updatesPerFile
}

func PatchReleasesFiles(releasesFile string, imageUpdates []ImageUpdate) error {
	var doc yaml.Node
	data, err := ioutil.ReadFile(releasesFile)
	if err != nil {
		return fmt.Errorf(`error reading %q: %w`, releasesFile, err)
	}
	err = yaml.Unmarshal(data, &doc)
	if err != nil {
		return fmt.Errorf(`error decoding yaml in %q: %w`, releasesFile, err)
	}
	releases := yamlNodeMapEntry(doc.Content[0], "releases")
	if releases.Kind != yaml.SequenceNode {
		return fmt.Errorf(`%s: "releases" is not a list`, releasesFile)
	}
	madeChanges := false
	for _, release := range releases.Content {
		name := yamlNodeMapEntry(release, "name")
		if name == nil || name.Kind != yaml.ScalarNode {
			continue
		}
		for _, update := range imageUpdates {
			if update.Release.Name != name.Value {
				continue
			}
			values := yamlNodeMapEntry(release, "values")
			if values == nil {
				continue
			}
			for _, chartValue := range values.Content {
				key := yamlNodeMapEntry(chartValue, "key")
				value := yamlNodeMapEntry(chartValue, "value")
				if key == nil || value == nil {
					continue
				}
				if key.Value == update.TagValue {
					value.Value = update.NewTag
					madeChanges = true
				}
			}
		}
	}
	if !madeChanges {
		return nil
	}
	return WriteIndentedYamlToFile(releasesFile, &doc)
}

func yamlNodeMapEntry(node *yaml.Node, name string) *yaml.Node {
	if node.Kind == yaml.MappingNode {
		for i := 0; i < len(node.Content); i += 2 {
			if node.Content[i].Kind == yaml.ScalarNode && node.Content[i].Value == name {
				return node.Content[i+1]
			}
		}
	}
	return nil
}
