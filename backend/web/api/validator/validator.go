// Input validation helper based on Alex Edwards' package from Let's Go Further
package validator

type Validator struct {
	errors map[string]string
}

func New() *Validator {
	v := Validator{
		errors: make(map[string]string),
	}
	return &v
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.errors[key]; !exists {
		v.errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func (v *Validator) Valid() bool {
	return len(v.errors) == 0
}

func (v *Validator) Errors() map[string]string {
	return v.errors
}
