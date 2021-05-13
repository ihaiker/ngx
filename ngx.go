package ngx

import (
	"github.com/ihaiker/ngx/v2/config"
	"github.com/ihaiker/ngx/v2/encoding"
	"github.com/ihaiker/ngx/v2/include"
	"github.com/ihaiker/ngx/v2/query"
)

var (
	MustParse     = config.MustParse
	MustParseWith = config.MustParseWith
	Parse         = config.Parse
	ParseWith     = config.ParseWith

	Selects = query.Selects
	Walk    = include.Walk

	Marshal            = encoding.Marshal
	MarshalWithOptions = encoding.MarshalWithOptions
	MarshalOptions     = encoding.MarshalOptions

	Unmarshal            = encoding.Unmarshal
	UnmarshalWithOptions = encoding.UnmarshalWithOptions
	UnmarshalDirectives  = encoding.UnmarshalDirectives
)

type Options = config.Options
type Configuration = config.Configuration
type Directive = config.Directive
type Directives = config.Directives
