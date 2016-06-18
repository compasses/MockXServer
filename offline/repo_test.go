package offline_test

import (
	"encoding/binary"
	"strconv"
	"testing"
)

func TestIntSliceEncode(t *testing.T) {
	var result []byte
	result = strconv.AppendInt(result, 124, 10)
	if string(result) != "124" {
		t.Errorf("%v\n", result)
	}
	result = strconv.AppendInt(result, 12, 10)

	another := make([]byte, 14)
	nul := binary.PutVarint(another, 2)
	nul = binary.PutVarint(another[4:], 111111111)
	if nul < 0 {
		t.Error()
	}
}
