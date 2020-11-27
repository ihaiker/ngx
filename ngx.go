package ngx

import (
	"github.com/ihaiker/ngx/config"
	"github.com/ihaiker/ngx/encoding"
	"github.com/ihaiker/ngx/include"
	"github.com/ihaiker/ngx/query"
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
