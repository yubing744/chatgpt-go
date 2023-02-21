package pkg

import "fmt"

type Log struct{}

func (log *Log) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
