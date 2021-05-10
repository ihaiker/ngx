package include

import (
	"github.com/ihaiker/ngx/config"
	"github.com/ihaiker/ngx/query"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"path/filepath"
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
}
`
const serverConfig = `
server  {
    # aginx api
    listen 80;
    server_name aginx.x.do;
    location / {
        proxy_pass http://127.0.0.1:8011;
        proxy_set_header X-Scheme $scheme;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'Upgrade';
    }
}
`

type includeAnalysis struct {
	suite.Suite
	includeFile string
	cfg         *config.Configuration
}

func (p *includeAnalysis) SetupTest() {
	p.includeFile = filepath.Join(os.TempDir(), "test.conf")
	err := ioutil.WriteFile(p.includeFile, []byte(serverConfig), 0644)
	p.Nil(err)
}

func (p *includeAnalysis) parse(options *config.Options) {
	var err error
	p.cfg, err = config.ParseWith([]byte(nginxConfig), options)
	p.Nil(err)

	err = Walk(p.cfg, func(d *config.Directive) bool {
		return d.Name == "include"
	}, func(args ...string) (files []string, err error) {
		files = []string{p.includeFile}
		return
	}, options)

	p.Nil(err)
}

func (p *includeAnalysis) TestNoMerge() {
	p.parse(&config.Options{MergeInclude: false})

	ds, err := query.Selects(p.cfg, "http", "include")
	p.Nil(err)
	p.Len(ds, 2)
	p.Equal(ds[0].Args[0], "mime.types")
	p.Equal(ds[1].Args[0], "conf.d/*.conf")
}

func (p *includeAnalysis) TestMerge() {
	p.parse(&config.Options{MergeInclude: true})

	ds, err := query.Selects(p.cfg, "http", "server")
	p.Nil(err)
	p.Len(ds, 1)
}

func TestInclude(t *testing.T) {
	suite.Run(t, new(includeAnalysis))
}
