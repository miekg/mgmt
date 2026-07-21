package main

import (
	"context"
	"flag"
	"log"
	"os"
	"sync"

	"github.com/purpleidea/mgmt/lang"
	"github.com/purpleidea/mgmt/lang/inputs"
	"github.com/purpleidea/mgmt/lang/interfaces"
	"github.com/purpleidea/mgmt/util"
	"github.com/spf13/afero"
)

// do a perf of initial langage loading.

var flagPerf = flag.Bool("p", true, "enable/disable perf gathering")

func main() {
	flag.Parse()

	for _, arg := range flag.Args() {
		data, err := os.ReadFile(arg)
		if err != nil {
			log.Fatal(err)
		}
		code := string(data)
		log.Printf("** %s", arg)

		// copied from lang/lang_test.go
		mmFs := afero.NewMemMapFs()
		afs := &afero.Afero{Fs: mmFs}
		fs := &util.AferoFs{Afero: afs}

		output, err := inputs.ParseInput(code, fs)
		if err != nil {
			log.Fatal(err)
		}
		for _, fn := range output.Workers {
			if err := fn(fs); err != nil {
				log.Fatal(err)
			}
		}

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		lang := &lang.Lang{
			Fs:    fs,
			Input: "/" + interfaces.MetadataFilename, // start path in fs
			Data: &lang.Data{
				UnificationStrategy: make(map[string]string),
			},
			Debug: false,
			Logf:  log.Printf,
		}

		if err := lang.Init(ctx); err != nil {
			log.Fatal(err)
		}

		wg := &sync.WaitGroup{}
		wg.Go(func() {
			if err := lang.Run(ctx); err != nil && err != context.Canceled {
				log.Fatal(err)
			}
		})
		cancel()
		wg.Wait()
		lang.Cleanup()
	}
}
