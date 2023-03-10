package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/m273d15/hcl2diff/pkg/hcl2diff"
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

func main() {
	var opts Hcl2DiffOpts
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	jsonPatch := hcl2diff.Hcl2DiffJsonWithGetter(opts.OldSrcVersion, opts.NewSrcVersion, opts.FileExtensions, opts.WorkDir)
	jsonB, err := json.Marshal(jsonPatch)
	if err != nil {
		log.Fatal("Failed to marshal json patch")
	}
	fmt.Println(string(jsonB))
}
