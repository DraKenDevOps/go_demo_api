package utils

import "fmt"

func InspectVariable(val any) {
	switch v := val.(type) {
	case string:
		fmt.Printf("It's a string: %q\n", v)
	case int:
		fmt.Printf("It's an integer: %d\n", v)
	case bool:
		fmt.Printf("It's a boolean: %t\n", v)
	default:
		fmt.Printf("Unknown type: %T\n", v) // %T prints the type name
	}
}
