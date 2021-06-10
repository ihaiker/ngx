package config

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

const nginxConfig = `
# test 
user nobody;
worker_processes auto;
events  {
    worker_connections 1024;
}
http  {
    include mime.types;
    default_type application/octet-stream;
    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
					'$status $body_bytes_sent "$http_referer" '
					'"$http_user_agent" "$http_x_forwarded_for"';
    # access_log  /var/log/nginx/access.log  main;
    sendfile on;
    # tcp_nopush     on;
    keepalive_timeout 65;
    gzip on;
    include conf.d/*.conf;
    include hosts.d/*.conf;
    include docker.d/*.conf;
}
`

type parseSuite struct {
	suite.Suite
}

func (p *parseSuite) TestParse() {
	config, err := ParseBytes([]byte(nginxConfig))
	p.Nil(err)
	p.Len(config.Body, 5)
	p.Equal(config.Body[0].Name, "#")

	user := config.Body.Get("user")
	p.NotNil(user)
	p.Len(user.Args, 1)
	p.Equal(user.Args[0], "nobody")
	p.Nil(user.Body)

	events := config.Body[3]
	p.Equal(events.Name, "events")
	p.Len(events.Body, 1)
	p.Equal(events.Body[0].Name, "worker_connections")
	p.Len(events.Body[0].Args, 1)
	p.Equal(events.Body[0].Args[0], "1024")

	p.T().Log(config.Pretty())
}

func (p parseSuite) TestNginxConfig() {
	conf, err := Parse("../query/_testdata/nginx.conf")
	p.Nil(err)
	p.T().Log(conf.Pretty())
}

func TestParseError(t *testing.T) {
	_, err := ParseBytes([]byte(`name test`))
	require.NotNil(t, err)
	t.Log(err)
}

func TestParse(t *testing.T) {
	suite.Run(t, new(parseSuite))
}
