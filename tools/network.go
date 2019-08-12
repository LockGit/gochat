/**
 * Created by lock
 * Date: 2019-08-12
 * Time: 16:00
 */
package tools

import (
	"fmt"
	"strings"
)

const (
	networkSplit = "@"
)

func ParseNetwork(str string) (network, addr string, err error) {
	if idx := strings.Index(str, networkSplit); idx == -1 {
		err = fmt.Errorf("addr: \"%s\" error, must be network@tcp:port or network@unixsocket", str)
		return
	} else {
		network = str[:idx]
		addr = str[idx+1:]
		return
	}
}
