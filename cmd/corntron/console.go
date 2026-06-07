package corntron

import (
	"corntron/internal"
)

func IsInTerminal() bool {
	return internal.IsInTerminal()
}

func Version() string {
	return internal.Version
}
