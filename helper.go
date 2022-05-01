package ipldswh

import (
	"bufio"
	"strconv"
)

func readTerminatedNumber(rd *bufio.Reader, delim byte) (int, error) {
	lstr, err := rd.ReadString(delim)
	if err != nil {
		return 0, err
	}
	lstr = lstr[:len(lstr)-1]

	return strconv.Atoi(lstr)
}
