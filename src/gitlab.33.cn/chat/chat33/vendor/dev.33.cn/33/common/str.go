package common

import (
	"bytes"
	"fmt"
)

func ToString(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

func StringConnect(data []string) string {
	buffer := bytes.Buffer{}

	for _, v := range data {
		_, err := buffer.WriteString(v)
		if err != nil {
			panic(err)
		}
	}

	return buffer.String()
}
