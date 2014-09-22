// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package validation

type Validation struct {
	errs map[string]*errors
}

func (v *Validation) HasErrors() bool {
	return len(v.errs) > 0
}

func (v *Validation) GetErrors() map[string]*errors {
	return v.errs
}

//
func (v *Validation) Apply(expr bool, msg, id string) *Validation {
	if expr {
		return v
	}

	if errs, found := v.errs[id]; found {
		errs.add(msg)
	} else {
		v.errs[id] = newErrors(msg)
	}

	return v
}
