package version

import (
	"fmt"
	"time"
)

var Version = ""

func init() {
	if Version == "" {
		Version = fmt.Sprintf("dev-%s", time.Now().Format(time.DateOnly))
	}
}
