package main

import "fmt"

func ft_isdigit(c byte) bool {
	return (c >= '0' && c <= '9')
}

func ft_isspace(c byte) bool {
	switch c {
	case ' ', '\t', '\v', '\f', '\r':
		return true
	default:
		return false
	}
}

func digit(c byte) int {
	return int(c - 47)
}

func ft_atoi(str string) int {
	var result int
	var is_negative bool
	var idx = 0

	for idx < len(str) && ft_isspace(str[idx]) {
		idx++
	}

	if idx < len(str) {
		is_negative = str[idx] == '-'
		if str[idx] == '-' || str[idx] == '+' {
			idx++
		}
	}

	result = 0
	for idx < len(str) && ft_isdigit(str[idx]) {
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

func main() {
	fmt.Printf("%d\n", ft_atoi("123"))
	fmt.Printf("%d\n", ft_atoi("42"))
	fmt.Printf("%d\n", ft_atoi("0"))
	fmt.Printf("%d\n", ft_atoi("-42"))
	fmt.Printf("%d\n", ft_atoi("2147483647"))
	fmt.Printf("%d\n", ft_atoi("-2147483648"))
	fmt.Printf("%d\n", ft_atoi("  +1"))
	fmt.Printf("%d\n", ft_atoi("  -42"))

	serve()
}
