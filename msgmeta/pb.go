package msgmeta

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"

	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
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

	if code == nil {
		return nil, errors.New("new buffer error")
	}

	err = binary.Write(code, binary.BigEndian, id)
	if err != nil {
		return nil, err
	}
	err = binary.Write(code, binary.BigEndian, rawdata)
	if err != nil {
		return nil, err
	}
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
