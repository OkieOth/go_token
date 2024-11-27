package token

import (
	"fmt"
	"time"

	"github.com/okieoth/goptional"
)

type TokenContent struct {
	Token            string
	ExirationSeconds int64
	LastUpdated      goptional.Optional[time.Time]
	LastChecked      goptional.Optional[time.Time]
}

func Dummy() {
	fmt.Printf(":-) %v/n", time.Now())
}
