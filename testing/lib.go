/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package testing

import (
	"os"
	"path/filepath"
	"runtime"

	"gomodules.xyz/encoding/yaml"
)

var rootDir = func() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Dir(filepath.Dir(b))
}()

func Load(filename string) (map[string]interface{}, error) {
	data, err := os.ReadFile(filepath.Join(rootDir, filename))
	if err != nil {
		return nil, err
	}

	var obj map[string]interface{}
	err = yaml.Unmarshal(data, &obj)
	return obj, err
}
