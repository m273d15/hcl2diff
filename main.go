package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/m273d15/hcl2diff/internal/getter"
	"github.com/m273d15/hcl2diff/pkg/hcl2json"
)

const (
	allowedCodeFormat = "Can be in format `git::https://github.com/m273d15/terraform-example?ref=main` or `file:///home/user/git/my-repo`"
)

type Opts struct {
	OldSrcVersion  string   `short:"o" long:"old-src" description:"Old hcl2 code path." required:"true"`
	NewSrcVersion  string   `short:"n" long:"new-src" description:"New hcl2 code path." required:"true"`
	FileExtensions []string `short:"e" long:"file-extension" default:".hcl" description:"File extension to search for to create a diff (defaults: '.hcl')." required:"false"`
	WorkDir        *string  `long:"workdir" description:"Workdir used to copy the source to (default: '$PWD/.temp_hcl2diff')" required:"false"`
}

func main() {
	var opts Opts
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	var wd string
	if opts.WorkDir == nil {
		wd = getter.InitWorkdir()
	} else {
		wd = getter.InitWorkdirPath(*opts.WorkDir)
	}
	// TODO: Handle getter issues
	firstHcl2Files := getter.GetFiles(opts.OldSrcVersion, wd, "oldSrc", opts.FileExtensions)
	secondHcl2Files := getter.GetFiles(opts.NewSrcVersion, wd, "newSrc", opts.FileExtensions)

	jsonPatch := hcl2json.Hcl2DiffJson(firstHcl2Files, secondHcl2Files)
	fmt.Println(string(jsonPatch))
}
