package chainsync

import (
	"github.com/cloudstruct/go-ouroboros-network/block"
	"github.com/cloudstruct/go-ouroboros-network/utils"
	"github.com/fxamacker/cbor/v2"
)

type WrappedBlock struct {
	// Tells the CBOR decoder to convert to/from a struct and a CBOR array
	_         struct{} `cbor:",toarray"`
	BlockType uint
	BlockCbor cbor.RawMessage
}

func NewWrappedBlock(blockType uint, blockCbor []byte) *WrappedBlock {
	return &WrappedBlock{
		BlockType: blockType,
		BlockCbor: blockCbor,
	}
}

type WrappedHeader struct {
	// Tells the CBOR decoder to convert to/from a struct and a CBOR array
	_          struct{} `cbor:",toarray"`
	Era        uint
	RawMessage cbor.RawMessage
	byronType  uint
	byronSize  uint
	headerCbor []byte
}

func NewWrappedHeader(era uint, byronType uint, blockCbor []byte) *WrappedHeader {
	w := &WrappedHeader{
		Era:       era,
		byronType: byronType,
	}
	// Record the original block size for Byron blocks
	if era == block.BLOCK_HEADER_TYPE_BYRON {
		// TODO: figure out why we have to add 2 to the length to match official message CBOR
		w.byronSize = uint(len(blockCbor)) + 2
	}
	// Parse block and extract header
	tmp := []cbor.RawMessage{}
	// TODO: figure out a better way to handle an error
	if err := cbor.Unmarshal(blockCbor, &tmp); err != nil {
		return nil
	}
	w.headerCbor = tmp[0]
	return w
}

func (w *WrappedHeader) UnmarshalCBOR(data []byte) error {
	var tmpHeader struct {
		// Tells the CBOR decoder to convert to/from a struct and a CBOR array
		_         struct{} `cbor:",toarray"`
		Era       uint
		HeaderRaw cbor.RawMessage
	}
	if err := cbor.Unmarshal(data, &tmpHeader); err != nil {
		return err
	}
	w.Era = tmpHeader.Era
	switch w.Era {
	case block.BLOCK_HEADER_TYPE_BYRON:
		var wrappedHeaderByron wrappedHeaderByron
		if _, err := utils.CborDecode(tmpHeader.HeaderRaw, &wrappedHeaderByron); err != nil {
			return err
		}
		w.byronType = wrappedHeaderByron.Metadata.Type
		w.byronSize = wrappedHeaderByron.Metadata.Size
		w.headerCbor = wrappedHeaderByron.RawHeader.Content.([]byte)
	default:
		var tag cbor.Tag
		if _, err := utils.CborDecode(tmpHeader.HeaderRaw, &tag); err != nil {
			return err
		}
		w.headerCbor = tag.Content.([]byte)
	}
	return nil
}

func (w *WrappedHeader) MarshalCBOR() ([]byte, error) {
	ret := []interface{}{
		w.Era,
	}
	switch w.Era {
	case block.BLOCK_HEADER_TYPE_BYRON:
		tmp := []interface{}{
			[]interface{}{
				w.byronType,
				w.byronSize,
			},
			cbor.Tag{
				Number:  24,
				Content: w.headerCbor,
			},
		}
		ret = append(ret, tmp)
	default:
		tag := cbor.Tag{
			Number:  24,
			Content: w.headerCbor,
		}
		ret = append(ret, tag)
	}
	cborData, err := utils.CborEncode(ret)
	if err != nil {
		return nil, err
	}
	return cborData, nil
}

func (w *WrappedHeader) HeaderCbor() []byte {
	return w.headerCbor
}

func (w *WrappedHeader) ByronType() uint {
	return w.byronType
}

type wrappedHeaderByron struct {
	// Tells the CBOR decoder to convert to/from a struct and a CBOR array
	_        struct{} `cbor:",toarray"`
	Metadata struct {
		// Tells the CBOR decoder to convert to/from a struct and a CBOR array
		_    struct{} `cbor:",toarray"`
		Type uint
		Size uint
	}
	RawHeader cbor.Tag
}
