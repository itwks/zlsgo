package zvalid

import (
	"testing"

	"github.com/sohaha/zlsgo"
)

func TestVar(t *testing.T) {
	tt := zlsgo.NewTest(t)

	var str string
	err := Var(&str, Text("is var").RemoveSpace())
	tt.EqualNil(err)
	tt.Equal("isvar", str)

	var i int
	err = Var(&i, Text("is var").RemoveSpace())
	tt.Equal(true, err != nil)
	tt.Equal(0, i)
	err = Var(&i, Text("99").RemoveSpace())
	tt.EqualNil(err)
	tt.Equal(99, i)

	var sts []string
	err = Var(&sts, Text("1,2,3,go").Separator(","))
	tt.EqualNil(err)
	tt.Equal([]string{"1", "2", "3", "go"}, sts)

	var data struct {
		Name string
	}

	err = Batch(
		BatchVar(&data.Name, Text("yes name")),
	)
	tt.EqualNil(err)
	tt.Equal("yes name", data.Name)

}

func TestVarDefault(t *testing.T) {
	tt := zlsgo.NewTest(t)

	var email string
	err := Var(&email, Text("email").IsMail())
	t.Log(email, err)
	tt.EqualExit(email, "")
	tt.EqualTrue(err != nil)

	err = Var(&email, Text("email").IsMail().Default("qq@qq.com"))
	t.Log(email, err)
	tt.EqualExit(email, "qq@qq.com")
	tt.EqualTrue(err != nil)

	var nu int
	err = Var(&nu, Text("Number").IsNumber().Default(123))
	t.Log(nu, err)
	tt.EqualTrue(err != nil)
	tt.EqualExit(nu, 123)

	var b bool
	err = Var(&b, Text("true").IsBool().Default(false))
	t.Log(b, err)
	tt.EqualTrue(err == nil)
	tt.EqualExit(b, true)
}
