package encoding_test

import (
	"github.com/ihaiker/ngx/encoding"
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

/*func (t *TestMarshalDemo) MarshalNgx() (config.Directives, error) {
	items := config.Directives{}
	items = append(items,config.New("name",t.Name))
	return items, nil
}*/

func TestMarshal(t *testing.T) {
	tm := new(TestMarshalDemo)
	tm.Name = "test"
	bs, err := encoding.Marshal(tm)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bs))
}

func TestArrays(t *testing.T) {
	names := []string{"n1", "n2"}
	bs, err := encoding.Marshal(names)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bs))
}

func TestArraysInt(t *testing.T) {
	names := []int{1, 2, 3, 4, 5}
	bs, err := encoding.Marshal(names)
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
	bs, err := encoding.Marshal(names)
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
	bs, err := encoding.Marshal(names)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bs))
}

func TestMarshal2(t *testing.T) {
	test := &Test{
		Name:    "姓名",
		Address: "地址",
		Create:  time.Now(),
		Age:     12,
		Sub: &TestSub{
			Demo: "实例",
		},
		Attr: map[string]string{
			"attr1": "a1",
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
	bs, err := encoding.Marshal(test)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bs))
}
