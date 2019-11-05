package updates

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

const indentLevel = 2

// WriteIndentedYamlToFile updates `fileName` atomically with indented yaml from `node`
func WriteIndentedYamlToFile(fileName string, node *yaml.Node) error {
	tmpFile, err := ioutil.TempFile(path.Dir(fileName), path.Base(fileName)+"*")
	if err != nil {
		return fmt.Errorf(`error creating tmpfile for %q: %w`, fileName, err)
	}
	//noinspection GoDeferInLoop
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	encoder := yaml.NewEncoder(tmpFile)
	encoder.SetIndent(indentLevel)
	if err = encoder.Encode(node); err != nil {
		return fmt.Errorf(`error re-encoding: %w`, err)
	}
	if err = os.Rename(tmpFile.Name(), fileName); err != nil {
		return fmt.Errorf(`error renaming %q to %q: %w`, tmpFile.Name(), fileName, err)
	}
	return nil
}
