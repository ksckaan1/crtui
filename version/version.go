package version

import (
	"fmt"
	"time"
)

var Version = fmt.Sprintf("dev-%s", time.Now().Format(time.DateOnly))
