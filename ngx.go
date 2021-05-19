package ngx

import (
	"github.com/ihaiker/ngx/v2/config"
	"github.com/ihaiker/ngx/v2/encoding"
)

var (
	MustParse     = config.MustParse
	MustParseIO   = config.MustParseIO
	MustParseWith = config.MustParseBytes

	Parse     = config.Parse
	ParseIO   = config.ParseIO
	ParseWith = config.ParseBytes

	Marshal            = encoding.Marshal
	MarshalWithOptions = encoding.MarshalWithOptions
	MarshalOptions     = encoding.MarshalOptions

	Unmarshal            = encoding.Unmarshal
	UnmarshalWithOptions = encoding.UnmarshalWithOptions
	UnmarshalDirectives  = encoding.UnmarshalDirectives
)

type Configuration = config.Configuration
type Directive = config.Directive
type Directives = config.Directives
