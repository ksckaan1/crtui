package ui

import (
	"time"

	"charm.land/bubbles/v2/spinner"
)

var LoadingSpinner = spinner.Spinner{
	Frames: []string{
		"\033[38;2;230;0;118ml\033[0m\033[37moading\033[0m",
		"\033[37ml\033[0m\033[38;2;230;0;118mo\033[0m\033[37mading\033[0m",
		"\033[37mlo\033[0m\033[38;2;230;0;118ma\033[0m\033[37mding\033[0m",
		"\033[37mloa\033[0m\033[38;2;230;0;118md\033[0m\033[37ming\033[0m",
		"\033[37mload\033[0m\033[38;2;230;0;118mi\033[0m\033[37mng\033[0m",
		"\033[37mloadi\033[0m\033[38;2;230;0;118mn\033[0m\033[37mg\033[0m",
		"\033[37mloadin\033[0m\033[38;2;230;0;118mg\033[0m",
	},
	FPS: time.Second / 7,
}
