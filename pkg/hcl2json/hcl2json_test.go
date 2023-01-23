package hcl2json

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func createTempFile(name string, t *testing.T) *os.File {
	baseDir := os.TempDir()

	file, err := os.Create(filepath.Join(baseDir, name))
	if err != nil {
		log.Println(err.Error())
		t.FailNow()
	}
	return file

}

func TestHcl2DiffJsonDeprecatedVariable(t *testing.T) {
	hclFile := hclwrite.NewEmptyFile()
	rootBody := hclFile.Body()

	variableBlock := rootBody.AppendNewBlock("variable", []string{"my_var"})
	variableBlockBody := variableBlock.Body()
	variableBlockBody.SetAttributeValue("description", cty.StringVal("This is my variable"))

	tfFileOld := createTempFile("variables_old.tf", t)
	tfFileOld.Write(hclFile.Bytes())

	variableBlockBody.SetAttributeValue("description", cty.StringVal("[Deprecated] Variable not needed anymore."))

	tfFileNew := createTempFile("variables_new.tf", t)
	tfFileNew.Write(hclFile.Bytes())

	expectedFrom := ""
	expectedPath := "/variable/my_var/0/description"
	expectedResult := JsonMap{
		"variable": JsonMap{
			"my_var": JsonMap{
				"description": JsonMap{
					"_operation": Op{
						Type:     "replace",
						From:     &expectedFrom,
						Path:     &expectedPath,
						OldValue: "[Deprecated] Variable not needed anymore.",
						Value:    "This is my variable",
					},
				},
			},
		},
	}
	diffJson := Hcl2DiffJsonMap([]string{tfFileNew.Name()}, []string{tfFileOld.Name()})

	assert.Equal(t, expectedResult, diffJson)
}
