package hooks

import "github.com/ihaiker/ngx/v2/config"

type Hook func(conf *config.Configuration) (err error)
