package formatter

import (
	"fmt"
	"strings"
)

func FormatBRL(value float64) string {
	// Convert to string with 2 decimal places
	str := fmt.Sprintf("%.2f", value)

	// Replace dot with comma
	str = strings.Replace(str, ".", ",", 1)

	// Add dots for thousands
	parts := strings.Split(str, ",")
	num := parts[0]
	dec := parts[1]

	var result []string
	for i := len(num) - 1; i >= 0; i-- {
		result = append(result, string(num[i]))
		if i > 0 && (len(num)-i)%3 == 0 {
			result = append(result, ".")
		}
	}

	// Reverse the slice and join everything
	for i := 0; i < len(result)/2; i++ {
		result[i], result[len(result)-1-i] = result[len(result)-1-i], result[i]
	}

	return fmt.Sprintf("R$ %s,%s", strings.Join(result, ""), dec)
}

func FormatUSD(value float64) string {
	// Convert to string with 2 decimal places
	str := fmt.Sprintf("%.2f", value)

	// Replace dot with comma
	str = strings.Replace(str, ".", ",", 1)

	// Add dots for thousands
	parts := strings.Split(str, ",")
	num := parts[0]
	dec := parts[1]

	var result []string
	for i := len(num) - 1; i >= 0; i-- {
		result = append(result, string(num[i]))
		if i > 0 && (len(num)-i)%3 == 0 {
			result = append(result, ".")
		}
	}

	// Reverse the slice and join everything
	for i := 0; i < len(result)/2; i++ {
		result[i], result[len(result)-1-i] = result[len(result)-1-i], result[i]
	}

	return fmt.Sprintf("US$ %s,%s", strings.Join(result, ""), dec)
}
