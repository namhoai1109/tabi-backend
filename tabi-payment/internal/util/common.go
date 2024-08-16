package util

// InterfaceToArrayString converts []interface{} to []string
func InterfaceToArrayString(in []interface{}) []string {
	out := []string{}
	for _, v := range in {
		out = append(out, v.(string))
	}
	return out
}

// TernaryOperator returns trueVal if condition is true, otherwise returns falseVal
func TernaryOperator(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}

	return falseVal
}
