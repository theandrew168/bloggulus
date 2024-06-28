package validator_test

import (
	"testing"

	"github.com/theandrew168/bloggulus/backend/test"
	"github.com/theandrew168/bloggulus/backend/web/validator"
)

func TestValidator(t *testing.T) {
	t.Parallel()

	v1 := validator.New()
	test.AssertEqual(t, v1.Valid(), true)

	v1.AddError("foo", "bar")
	test.AssertEqual(t, v1.Valid(), false)
	test.AssertEqual(t, v1.Errors(), map[string]string{"foo": "bar"})

	v2 := validator.New()
	test.AssertEqual(t, v2.Valid(), true)

	v2.Check(false, "foo", "bar")
	test.AssertEqual(t, v2.Valid(), false)
	test.AssertEqual(t, v2.Errors(), map[string]string{"foo": "bar"})
}
