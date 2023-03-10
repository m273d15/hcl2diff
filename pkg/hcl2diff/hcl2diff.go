package hcl2diff

import (
	"github.com/m273d15/hcl2diff/internal/getter"
	"github.com/m273d15/hcl2diff/pkg/hcl2json"
)

const (
	allowedCodeFormat = "Can be in format `git::https://github.com/m273d15/terraform-example?ref=main` or `file:///home/user/git/my-repo`"
)

type Hcl2DiffOpts struct {
	OldSrcVersion  string   `short:"o" long:"old-src" description:"Old hcl2 code path." required:"true"`
	NewSrcVersion  string   `short:"n" long:"new-src" description:"New hcl2 code path." required:"true"`
	FileExtensions []string `short:"e" long:"file-extension" default:".hcl" description:"File extension to search for to create a diff (defaults: '.hcl')." required:"false"`
	WorkDir        *string  `long:"workdir" description:"Workdir used to copy the source to (default: '$PWD/.temp_hcl2diff')" required:"false"`
}

func Hcl2DiffJsonWithGetter(oldSrcVersion, newSrcVersion string, fileExtensions []string, workDir *string) hcl2json.JsonMap {
	var wd string
	if workDir == nil {
		wd = getter.InitWorkdir()
	} else {
		wd = getter.InitWorkdirPath(*workDir)
	}
	// TODO: Handle getter issues
	oldHcl2Files := getter.GetFiles(oldSrcVersion, wd, "oldSrc", fileExtensions)
	newHcl2Files := getter.GetFiles(newSrcVersion, wd, "newSrc", fileExtensions)
	return hcl2json.Hcl2DiffJsonMap(oldHcl2Files, newHcl2Files)
}
