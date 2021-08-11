package testing

import (
	"io/ioutil"
	"path/filepath"
	"runtime"

	"k8s.io/apimachinery/pkg/util/yaml"
)

var rootDir = func() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Dir(filepath.Dir(b))
}()

func Load(filename string) (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(filepath.Join(rootDir, filename))
	if err != nil {
		return nil, err
	}

	var obj map[string]interface{}
	err = yaml.Unmarshal(data, &obj)
	return obj, err
}
