package encoding

import (
	"github.com/ihaiker/ngx/config"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

type Test struct {
	Name    string     `ngx:"name"`
	Address string     `ngx:"address"`
	Create  *time.Time `ngx:"time,2006-01-02 15:04:05"`
	Age     int

	Sub *TestSub `ngx:"sub"`

	Attr  map[string]string `ngx:"attr"`
	Demos map[string]*TestSub

	Tags []string `ngx:"tags"`

	DemoAry []*TestSub `ngx:"demoAry"`
}

type TestSub struct {
	Demo string `ngx:"demo"`
	Attr string `ngx:"attr"`
}

func (t *TestSub) UnmarshalNgx(items config.Directives) error {
	for _, item := range items {
		if d := item.Body.Get("demo"); d != nil {
			t.Demo = strings.Join(d.Args, "")
		}
	}
	return nil
}

func TestUnmarshal(t *testing.T) {
	tt := new(Test)
	data := []byte(`
		name: "姓名";
		address: "地址";
		time: "2020-04-04 04:04:04";
		Age: 20;
        sub {
			demo: "测试模板";
		}
		attr: name "zhou";
		attr: age 12312;
		attr {
			address: "地址一样";
			port: 1024;
		}
		Demos d1 {
			demo: "d1的值";
		}
		Demos {
			d2 {
				demo: "d2的值，";
			}
			d3 {
				demo: "d3的值";
			}
		}
		tags: t1 t2 t3 "t1 t2";
		demoAry {
			demo: "demo ary 1";
		}
		demoAry {
			demo: "demo ary 2";
		}
	`)

	err := Unmarshal(data, tt)
	require.Nil(t, err)
	require.Equal(t, tt.Name, "姓名")
	require.Equal(t, tt.Address, "地址")
	require.Len(t, tt.Demos, 3)
	require.Equal(t, tt.Age, 20)
	require.Len(t, tt.Tags, 4)
	require.Equal(t, tt.Tags[1], "t2")
}
