package util_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/web/util"
)

func TestErrorsAdd(t *testing.T) {
	t.Parallel()

	e := util.NewErrors()
	test.AssertEqual(t, e.Valid(), true)

	e.Add("message")
	test.AssertEqual(t, e.Valid(), false)
}

func TestErrorsAddField(t *testing.T) {
	t.Parallel()

	e := util.NewErrors()
	test.AssertEqual(t, e.Valid(), true)

	e.AddField("message", "field")
	test.AssertEqual(t, e.Valid(), false)
}

func TestErrorsCheck(t *testing.T) {
	t.Parallel()

	e := util.NewErrors()
	test.AssertEqual(t, e.Valid(), true)

	e.Check(true, "message")
	test.AssertEqual(t, e.Valid(), true)

	e.Check(false, "message")
	test.AssertEqual(t, e.Valid(), false)
}

func TestErrorsCheckField(t *testing.T) {
	t.Parallel()

	e := util.NewErrors()
	test.AssertEqual(t, e.Valid(), true)

	e.CheckField(true, "message", "field")
	test.AssertEqual(t, e.Valid(), true)

	e.CheckField(false, "message", "field")
	test.AssertEqual(t, e.Valid(), false)
}
