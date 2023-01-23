package hcl2json

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/tmccombs/hcl2json/convert"
	"github.com/wI2L/jsondiff"
)

type JsonMap = map[string]interface{}

type Op struct {
	Type     string      `json:"op"`
	From     *string     `json:"from"`
	Path     *string     `json:"path"`
	OldValue interface{} `json:"oldValue"`
	Value    interface{} `json:"value"`
}

func hcl2Json(hcl2FilePath string, target *JsonMap) {
	hcl2File, err := ioutil.ReadFile(hcl2FilePath)
	if err != nil {
		log.Fatal("Cannot read the hcl file:", hcl2FilePath, err)
	}

	hcl2Json, err := convert.Bytes(hcl2File, "", convert.Options{})
	if err != nil {
		log.Fatal("Cannot convert hcl to JSON:", err)
	}
	json.Unmarshal(hcl2Json, &target)
}

// Create nested JSON structure which is helpful if used in OPA
func toNestedJson(patch jsondiff.Patch) JsonMap {
	jsonMap := make(JsonMap)
	for _, op := range patch {
		path := string(op.Path)
		pathElements := strings.Split(path, "/")[1:]
		var m = jsonMap
		for _, e := range pathElements {
			if _, err := strconv.Atoi(e); err == nil {
				// Skip index integers in path
				continue
			}

			if val, ok := m[e]; ok {
				m = val.(JsonMap)
			} else {
				m[e] = make(JsonMap)
				m = m[e].(JsonMap)
			}
		}
		mop := Op{
			Type:     op.Type,
			From:     (*string)(&op.From),
			Path:     (*string)(&op.Path),
			OldValue: op.OldValue,
			Value:    op.Value,
		}
		m["_operation"] = mop
	}
	return jsonMap
}

func jsonFromHcl2Files(files []string) JsonMap {
	var merged JsonMap
	for _, file := range files {
		hcl2Json(file, &merged)
	}
	return merged
}

func Hcl2DiffPatch(oldHcl2Files, newHcl2Files []string) jsondiff.Patch {
	mergedOld := jsonFromHcl2Files(oldHcl2Files)
	mergedNew := jsonFromHcl2Files(newHcl2Files)

	patch, err := jsondiff.Compare(mergedOld, mergedNew)
	if err != nil {
		log.Fatal("Cannot create a diff:", err)
	}
	return patch
}

func Hcl2DiffJsonMap(oldHcl2Files, newHcl2Files []string) JsonMap {
	patch := Hcl2DiffPatch(oldHcl2Files, newHcl2Files)
	jsonPatch := toNestedJson(patch)
	return jsonPatch
}

func Hcl2DiffJson(oldHcl2Files, newHcl2Files []string) []byte {
	jsonPatch := Hcl2DiffJsonMap(oldHcl2Files, newHcl2Files)
	b, err := json.Marshal(jsonPatch)
	if err != nil {
		log.Fatal("Failed to marshal json patch")
	}
	return b
}
