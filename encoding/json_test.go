package encoding

import (
	"github.com/ihaiker/ngx/v2/config"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

func TestJson(t *testing.T) {
	data, err := ioutil.ReadFile("../query/_testdata/nginx.conf")
	require.Nil(t, err)

	conf, err := config.ParseBytes(data)
	require.Nil(t, err)

	bs, err := JsonIndent(conf, "    ", "    ")
	require.Nil(t, err)

	t.Log(string(bs))
}
