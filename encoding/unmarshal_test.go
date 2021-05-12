package encoding

import (
	"github.com/ihaiker/ngx/config"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
	"time"
)

type Test struct {
	Name     string     `ngx:"name"`
	StarName *string    `ngx:"star_name"`
	Address  string     `ngx:"address"`
	Create   *time.Time `ngx:"time,2006-01-02 15:04:05"`
	Update   time.Time  `ngx:"time,2006-01-02 15:04:05"`

	Age int

	Sub  *TestSub `ngx:"sub"`
	Sub1 TestSub  `ngx:"sub1"`

	Sub2 *TestSubUnmarshaler `ngx:"sub2"`
	Sub3 TestSubUnmarshaler  `ngx:"sub3"`

	Attr  map[string]string `ngx:"attr"`
	Demos map[string]*TestSub

	Tags []string `ngx:"tags"`

	DemoAry []*TestSubUnmarshaler `ngx:"demoAry"`
}

type TestSub struct {
	Demo string `ngx:"demo"`
	Attr string `ngx:"attr"`
}

type TestSubUnmarshaler struct {
	Demo string `ngx:"demo"`
	Attr string `ngx:"attr"`
}

func (t *TestSubUnmarshaler) UnmarshalNgx(items *config.Configuration) error {
	if d := items.Body.Get("demo"); d != nil {
		t.Demo = strings.Join(d.Args, "")
	}
	t.Attr = "Unmarshal"
	return nil
}

type TestSubHandler struct {
}

func (t *TestSubHandler) MarshalNgx(v interface{}) (cfg *config.Configuration, err error) {
	ts := v.(*TestSub)
	cfg = config.Config()
	cfg.Body.Append(config.New("demo", ts.Demo))
	cfg.Body.Append(config.New("attr", ts.Attr))
	return
}

func (t *TestSubHandler) UnmarshalNgx(item *config.Configuration) (v interface{}, err error) {
	v = new(TestSub)

	if d := item.Body.Get("demo"); d != nil {
		v.(*TestSub).Demo = strings.Join(d.Args, "") + "@ngx"
	}
	if d := item.Body.Get("attr"); d != nil {
		v.(*TestSub).Attr = strings.Join(d.Args, "") + "@ngx"
	}
	return
}

type TestUnmarshalSuite struct {
	suite.Suite
	tt *Test
}

func (p *TestUnmarshalSuite) SetupTest() {
	p.tt = new(Test)
}

func (p *TestUnmarshalSuite) Unmarshal(ngx string) {
	err := Unmarshal([]byte(ngx), p.tt)
	p.Nil(err)
}

func (p *TestUnmarshalSuite) TestBase() {
	p.Unmarshal(`
		name: "姓名";
		address: "地址";
		time: "2020-04-04 04:04:04";
		Age: 20;
	`)
	p.Equal("姓名", p.tt.Name)
	p.Equal("地址", p.tt.Address)
	p.Equal(time.Date(2020, 4, 4, 4, 4, 4, 0, time.UTC).Unix(), p.tt.Create.Unix())
}

func (p *TestUnmarshalSuite) TestBaseMap() {
	p.Unmarshal(`
		attr: name "zhou";
		attr: age 12312;
		attr {
			address: "地址一样";
			port: 1024;
		}
	`)
	p.Len(p.tt.Attr, 4)
	p.Equal("zhou", p.tt.Attr["name"])
	p.Equal("12312", p.tt.Attr["age"])
	p.Equal("地址一样", p.tt.Attr["address"])
	p.Equal("1024", p.tt.Attr["port"])
}

func (p *TestUnmarshalSuite) TestStarSubStruct() {
	p.Unmarshal(`
		name: "star_sub_demo";
		sub {
			demo: "demo";
			attr: "attr";
		}
	`)
	p.Equal("star_sub_demo", p.tt.Name)
	p.NotNil(p.tt.Sub)
	p.Nil(p.tt.Sub2)
	p.Equal("demo", p.tt.Sub.Demo)
	p.Equal("attr", p.tt.Sub.Attr)
}

func (p *TestUnmarshalSuite) TestSubStruct() {
	p.Unmarshal(`
		name: "sub_demo";
		sub1 {
			demo: "demo";
			attr: "attr";
		}
	`)
	p.Equal("sub_demo", p.tt.Name)
	p.Nil(p.tt.Sub)
	p.Equal("demo", p.tt.Sub1.Demo)
	p.Equal("attr", p.tt.Sub1.Attr)
}

func (p *TestUnmarshalSuite) TestUnmarshaler() {
	{
		t := new(TestSubUnmarshaler)
		err := Unmarshal([]byte(``), t)
		p.Nil(err)
		p.Equal("Unmarshal", t.Attr)
	}
	{
		t := new(TestSubUnmarshaler)
		err := Unmarshal([]byte(`
			demo: "unmarshal demo";
		`), t)
		p.Nil(err)
		p.Equal("unmarshal demo", t.Demo)
		p.Equal("Unmarshal", t.Attr)
	}
}

func (p *TestUnmarshalSuite) TestSubUnmarshaler() {
	p.Unmarshal(`
		name: "sub_unmarshal";
		sub2 {
			demo: "sub2_demo";
			attr: "whatever";
		}
		sub3 {
			demo: "sub3_demo";
			attr: "whatever";
		}
	`)
	p.Equal("sub_unmarshal", p.tt.Name)

	p.NotNil(p.tt.Sub2)
	p.Equal("sub2_demo", p.tt.Sub2.Demo)
	p.Equal("Unmarshal", p.tt.Sub2.Attr)

	p.Equal("sub3_demo", p.tt.Sub3.Demo)
	p.Equal("Unmarshal", p.tt.Sub3.Attr)

}

func (p *TestUnmarshalSuite) TestMapStruct() {
	p.Unmarshal(`
		name: "test map[base]struct";
		Demos d1 {
			demo: "d1_demo_value";
			attr: "attr_value";
		}
		Demos {
			d2 {
				demo: "d2_demo";
			}
			d3 {
				demo: "d3_demo";
			}
		}
	`)
	p.Equal("test map[base]struct", p.tt.Name)
	p.Len(p.tt.Demos, 3)

	p.Contains(p.tt.Demos, "d1")
	p.Contains(p.tt.Demos, "d2")
	p.Contains(p.tt.Demos, "d3")

	p.NotNil(p.tt.Demos["d1"])
	p.NotNil(p.tt.Demos["d2"])
	p.NotNil(p.tt.Demos["d3"])

	p.Equal("d1_demo_value", p.tt.Demos["d1"].Demo)
	p.Equal("attr_value", p.tt.Demos["d1"].Attr)
	p.Equal("d2_demo", p.tt.Demos["d2"].Demo)
	p.Equal("d3_demo", p.tt.Demos["d3"].Demo)
}

func (p *TestUnmarshalSuite) TestBaseSlice() {
	p.Unmarshal(`
		name: "slice test";
		tags: t1 t2 t3 "t1 t3";
	`)
	p.Equal("slice test", p.tt.Name)
	p.Len(p.tt.Tags, 4)
	p.Equal("t2", p.tt.Tags[1])
	p.Equal("t1 t3", p.tt.Tags[3])
}

func (p *TestUnmarshalSuite) TestStructSlice() {
	p.Unmarshal(`
		demoAry {
			demo: "demo ary 1";
			attr: "attr1";
		}
		demoAry {
			demo: "demo ary 2";
		}
	`)
	p.Len(p.tt.DemoAry, 2)
	p.Equal("demo ary 1", p.tt.DemoAry[0].Demo)
	//the struct 实现 Unmarshaler
	p.Equal("Unmarshal", p.tt.DemoAry[0].Attr)

	p.Equal("demo ary 2", p.tt.DemoAry[1].Demo)
}

func (p TestUnmarshalSuite) TestTypeHandler() {
	RegTypeHandler(new(TestSub), new(TestSubHandler))

	ts := new(TestSub)
	err := Unmarshal([]byte(`
		demo: "demo";
		attr: "attr";
	`), ts)
	p.Nil(err)
	p.Equal("demo@ngx", ts.Demo)
	p.Equal("attr@ngx", ts.Attr)
}

func (p TestUnmarshalSuite) TestSubTypeHandler() {
	RegTypeHandler(new(TestSub), new(TestSubHandler))

	p.Unmarshal(`
		name: "type handler test";
		sub {
			demo: "demo";
			attr: "attr";
		}
	`)

	p.Equal("type handler test", p.tt.Name)
	p.NotNil(p.tt.Sub)
	p.Nil(p.tt.Sub2)
	p.Equal("demo@ngx", p.tt.Sub.Demo)
	p.Equal("attr@ngx", p.tt.Sub.Attr)
}

func TestUnmarshal(t *testing.T) {
	suite.Run(t, new(TestUnmarshalSuite))
}
