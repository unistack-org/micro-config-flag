package flag

import (
	"strings"
)

func getFlagOpts(tf string) (string, string, string) {
	ret := make([]string, 3)
	vals := strings.Split(tf, ",")
	f := 0
	for _, val := range vals {
		p := strings.Split(val, "=")
		switch p[0] {
		case "name":
			f = 0
		case "desc":
			f = 1
		case "default":
			f = 2
		default:
			ret[f] += "," + val
			continue
		}
		ret[f] = p[1]
	}
	for idx := range ret {
		if ret[idx][0] == '\'' {
			ret[idx] = ret[idx][1:]
		}
		if ret[idx][len(ret[idx])-1] == '\'' {
			ret[idx] = ret[idx][:len(ret[idx])-1]
		}
	}
	return ret[0], ret[1], ret[2]
}
