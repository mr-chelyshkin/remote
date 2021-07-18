package hooks

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"strings"
)

func ValidatePath(_path string) error {
	for _, blockDir := range viper.GetStringSlice("block_dirs_list") {
		if strings.HasPrefix(_path, blockDir) {
			return errors.Errorf("validate path error: path '%s' blocked on server side", _path)
		}
	}
	return nil
}
