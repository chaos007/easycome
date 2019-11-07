package msgmeta

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"

	log "github.com/sirupsen/logrus"
	"github.com/golang/protobuf/proto"
)

// BuildPacket 建包
func BuildPacket(data proto.Message) ([]byte, error) {
	id, ok := MessageMetaByName(proto.MessageName(data))
	if !ok {
		log.Errorln("can not find msgid return,msgid:", id)
		return nil, errors.New("cannot find message id")
	}
	rawdata, err := proto.Marshal(data)
	if err != nil {
		log.Errorf("Marshal:%s", err.Error())
		return nil, err
	}

	code := bytes.NewBuffer([]byte{})

	binary.Write(code, binary.BigEndian, id)
	binary.Write(code, binary.BigEndian, rawdata)

	return code.Bytes(), nil
}

// ParsePacket 解析包
func ParsePacket(msgID int32, pkt []byte) (data proto.Message, err error) {
	if name, ok := msgIDToName[msgID]; ok {
		data = reflect.New(proto.MessageType(name).Elem()).Interface().(proto.Message)
		err = proto.Unmarshal(pkt, data)
	}
	return
}
