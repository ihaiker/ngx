package hooks

import (
	"github.com/ihaiker/ngx/v2/config"
)

func SWitch(params *Parameters) Hook {
	return func(conf *config.Configuration) error {
		return setValue(conf.Body)
	}
}

func setValue(items config.Directives) error {
	return nil
}
