package ui

import (
	"fmt"

	"github.com/samber/lo"
)

func CountKind(count int, kindSingular, kindPlural string) string {
	return fmt.Sprintf("%d %s", count, lo.Ternary(count > 1, kindPlural, kindSingular))
}
