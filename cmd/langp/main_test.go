package main

import (
	"context"
	"io"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/purpleidea/mgmt/lang"
	"github.com/purpleidea/mgmt/lang/inputs"
	"github.com/purpleidea/mgmt/lang/interfaces"
	"github.com/purpleidea/mgmt/util"
	"github.com/spf13/afero"
)

var files = []string{
	"../../examples/lang/http-server0.mcl",
}

func BenchmarkLang(b *testing.B) {
	for _, arg := range files {
		data, err := os.ReadFile(arg)
		if err != nil {
			b.Fatal(err)
		}
		code := string(data)
		log.SetOutput(io.Discard)
		b.Run(arg, func(b *testing.B) {
			// copied from lang/lang_test.go
			mmFs := afero.NewMemMapFs()
			afs := &afero.Afero{Fs: mmFs}
			fs := &util.AferoFs{Afero: afs}

			b.ResetTimer()

			for b.Loop() {
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
					Input: "/" + interfaces.MetadataFilename,
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

		})
	}
}
