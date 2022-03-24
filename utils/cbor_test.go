package utils_test

import (
	"reflect"
	"testing"

	"github.com/cloudstruct/go-ouroboros-network/protocol"
	"github.com/cloudstruct/go-ouroboros-network/protocol/chainsync"
	"github.com/cloudstruct/go-ouroboros-network/utils"
)

var testArray = []byte("\x81\x00")

func TestCborEncode(t *testing.T) {
	msg := &chainsync.MsgRequestNext{
		MessageBase: protocol.MessageBase{
			MessageType: chainsync.MESSAGE_TYPE_REQUEST_NEXT,
		},
	}

	data, err := utils.CborEncode(msg)
	if err != nil {
		t.Errorf("Output %q not equal to expected %q.  Received error '%s'.", data, testArray, err)
	}
	if !reflect.DeepEqual(data, testArray) {
		t.Errorf("Output %q not equal to expected %q.", data, testArray)
	}
}

//func TestCborDecode(t *testing.T) {
//}
