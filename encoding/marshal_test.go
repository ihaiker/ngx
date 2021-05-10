package encoding

import (
	"fmt"
	"github.com/ihaiker/ngx/config"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

type (
	TestMarshalSubDemo struct {
		Address string `ngx:"address"`
	}
	TestMarshalDemo struct {
		Name       string              `ngx:"name"`
		SubAddress *TestMarshalSubDemo `ngx:"sub"`
	}
)

var test = &Test{
	Name:    "姓名",
	Address: "地址",
	Create:  new(time.Time),
	Age:     12,
	Sub: &TestSub{
		Demo: "实例",
	},
	Attr: map[string]string{
		"attr1": "a1 b2",
		"attr2": "a2",
	},
	Demos: map[string]*TestSub{
		"d1": &TestSub{
			Demo: "d1-1",
		},
		"d2": &TestSub{
			Demo: "d2---2",
		},
	},
	Tags: []string{"t1", "t2", "t3"},
	DemoAry: []*TestSub{
		{
			Demo: "t1-d2",
			Attr: "t1-attr",
		},
		{
			Demo: "t2-d3",
			Attr: "t2-attr",
		},
	},
}

func TestDate(t *testing.T) {
	d := time.Now()
	v, err := MarshalWithOptions(d, Options{
		DateFormat: "2006-01-02 15:04:05",
	})
	require.Nil(t, err)
	t.Log(string(v))
}

func TestMarshal(t *testing.T) {
	tm := new(TestMarshalDemo)
	tm.Name = "test"
	bs, err := Marshal(tm)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bs))
}

func TestArrays(t *testing.T) {
	names := []string{"n1", "n2"}
	bs, err := Marshal(names)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bs))
}

func TestArraysInt(t *testing.T) {
	names := []int{1, 2, 3, 4, 5}
	bs, err := Marshal(names)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bs))
}

func TestArrayStruct(t *testing.T) {
	names := []TestMarshalDemo{
		{
			Name: "name",
		},
	}
	bs, err := Marshal(names)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bs))
}

func TestStruct(t *testing.T) {
	names := &TestMarshalDemo{
		Name: "name",
		SubAddress: &TestMarshalSubDemo{
			Address: "地址",
		},
	}
	bs, err := Marshal(names)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bs))
}

func TestMarshal2(t *testing.T) {
	bs, err := Marshal(test)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bs))
}

type DateReg struct {
	Format string
}

func (d *DateReg) MarshalNgx(v interface{}) (config.Directives, error) {
	return config.Directives{
		config.New("date", strconv.Quote(v.(time.Time).Format(d.Format))),
	}, nil
}

func (d *DateReg) UnmarshalNgx(item config.Directives) (interface{}, error) {
	if len(item[0].Args) == 0 {
		return nil, fmt.Errorf("%s must has a param", item[0].Name)
	}
	return time.Parse(d.Format, item[0].Args[0])
}

func TestReg(t *testing.T) {
	tt := time.Now()
	test.Create = &tt

	Defaults.TypeHandlers.Reg(time.Now(), &DateReg{Format: "2006-01-02"})

	bs, err := MarshalWithOptions(test, *Defaults)
	require.Nil(t, err)
	t.Log(string(bs))
}
