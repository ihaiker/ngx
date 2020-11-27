package ngx

import (
	"github.com/xhaiker/ngx/config"
	"github.com/xhaiker/ngx/encoding"
	"github.com/xhaiker/ngx/include"
	"github.com/xhaiker/ngx/query"
)

var (
	MustParse     = config.MustParse
	MustParseWith = config.MustParseWith
	Parse         = config.Parse
	ParseWith     = config.ParseWith

	Selects = query.Selects
	Walk    = include.Walk

	Marshal = encoding.Marshal

	Unmarshal   = encoding.Unmarshal
	MarshalWith = encoding.UnmarshalWith
)

type Options = config.Options
