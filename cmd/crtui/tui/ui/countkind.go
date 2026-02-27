package ui

import (
	"fmt"

	"github.com/samber/lo"
)

func CountKind(count int, singular, plural string) string {
	return fmt.Sprintf("%d %s", count, lo.Ternary(count > 1, plural, singular))
}
