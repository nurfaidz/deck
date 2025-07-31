package helpers

import "fmt"

func FormatCurrency(amount uint) string {
	return fmt.Sprintf("Rp %d", amount)
}
