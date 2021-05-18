package method

import (
	"github.com/ihaiker/ngx/v2/config"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestArgEqual(t *testing.T) {
	items := config.Directives{
		{
			Name: "proxy_set_header",
			Args: []string{"Host", "$host"},
		},
		{
			Name: "proxy_set_header",
			Args: []string{"Upgrade", "$http_upgrade;"},
		},
	}
	items = arg(items, 0, "Host")
	require.Len(t, items, 1)
	require.Equal(t, "$host", items[0].Args[1])

	//arg(1 eq '')
}
