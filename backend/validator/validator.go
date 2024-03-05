// Input validation helper based on Alex Edwards' package from Let's Go Further
package validator

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	v := Validator{
		Errors: make(map[string]string),
	}
	return &v
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}
