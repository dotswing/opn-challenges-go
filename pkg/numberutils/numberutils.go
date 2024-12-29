package numberutils

import (
	"fmt"

	"github.com/dustin/go-humanize"
)

func FormatFloat(number float64) string {
	integerPart := int64(number)
	decimalPart := number - float64(integerPart)
	integerPartWithCommas := humanize.Comma(integerPart)
	return fmt.Sprintf("%s%s", integerPartWithCommas, fmt.Sprintf("%.2f", decimalPart)[1:])
}
