package atoi

func isdigit(c byte) bool {
	return (c >= '0' && c <= '9')
}

func isspace(c byte) bool {
	switch c {
	case ' ', '\t', '\v', '\f', '\r':
		return true
	default:
		return false
	}
}

func digit(c byte) int {
	return int(c - 48)
}

func Atoi(str string) int {
	var result int
	var is_negative bool
	var idx = 0

	for idx < len(str) && isspace(str[idx]) {
		idx++
	}

	if idx < len(str) {
		is_negative = str[idx] == '-'
		if str[idx] == '-' || str[idx] == '+' {
			idx++
		}
	}

	result = 0
	for idx < len(str) && isdigit(str[idx]) {
		result *= 10
		result -= digit(str[idx])
		idx++
	}

	if is_negative {
		return result
	} else {
		return -result
	}
}
