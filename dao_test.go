package riverlet

import "testing"

type AExample struct {
	Name string `json:"name"`
}

func (a AExample) Identity() string {
	return "a.example"
}

func TestBuiltinPersistence(t *testing.T) {
	p := NewBuiltinPersistence[AExample]()
	t.Log(p)
	var j = AExample{}
	j.Name = "dao_test"
	if o, e := p.Insert(j, 5, 20); nil != e {
		t.Error("Insert Failed:", e)
	} else {
		t.Log(">>>", o)
		if jj, ee := p.Get(o.Seq); nil != ee {
			t.Error("Get Failed:", ee)
		} else {
			t.Log(">>>", jj)
		}
	}
}
