package config

import (
	"fmt"
	"path/filepath"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
)

func Load(path string, cfg interface{}) error {
	if err := aconfig.LoaderFor(cfg, aconfig.Config{
		Files: []string{
			filepath.Join(path),
		},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yml": aconfigyaml.New(),
		},
	}).Load(); err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}
	return nil
}
