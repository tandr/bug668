package bug668

import (
	"log"
	"testing"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/proto"
)

func TestDbBlock_packUnpack(t *testing.T) {

	b1 := Block{
		BlockId:   "someid",
		Timestamp: ptypes.TimestampNow(),
	}

	packed, err := packBlock(&b1)
	if err != nil {
		t.Fatal(err)
	}

	b2, err := unpackBlock(packed)

	if err != nil {
		t.Fatal(err)
	}

	if b2.BlockId != b1.BlockId || !proto.Equal(b1.Timestamp, b2.Timestamp) {
		t.Fatal("blocks are different", b1, b2)
	}
}

// block wrapper as possible additional meta info holder
type dbBlock struct {
	*Block `protobuf:"bytes,1,opt,name=block" json:"block,omitempty"`
}

// fulfilling Message interface
func (m *dbBlock) Reset()         { *m = dbBlock{} }
func (m *dbBlock) String() string { return proto.CompactTextString(m) }
func (*dbBlock) ProtoMessage()    {}

func init() {
	proto.RegisterType((*dbBlock)(nil), "bug668.dbBlock")
}

func unpackBlock(value []byte) (*Block, error) {
	var dbBlock dbBlock
	err := proto.Unmarshal(value, &dbBlock)

	if err != nil {
		log.Fatalf("cannot un-marshall block %v, error=%v", dbBlock.BlockId, err)
	}

	return dbBlock.Block, err
}

func packBlock(b *Block) ([]byte, error) {
	dbBlock := dbBlock{b}

	buf, err := proto.Marshal(&dbBlock)

	if err != nil {
		log.Fatalf("cannot marshall block %v, error=%v", dbBlock.BlockId, err)
	}

	return buf, err
}
