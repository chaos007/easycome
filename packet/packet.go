package packet

import (
	"errors"
	"math"
)

// 包限制长度
const (
	PacketLimit = 65535
)

// Packet 数据包
type Packet struct {
	pos  int
	data []byte
}

// Data data
func (p *Packet) Data() []byte {
	return p.data
}

// Length 长度
func (p *Packet) Length() int {
	return len(p.data)
}

// ReadBool 读bool
func (p *Packet) ReadBool() (ret bool, err error) {
	b, _err := p.ReadByte()

	if b != byte(1) {
		return false, _err
	}

	return true, _err
}

// ReadByte 读字节
func (p *Packet) ReadByte() (ret byte, err error) {
	if p.pos >= len(p.data) {
		err = errors.New("read byte failed")
		return
	}

	ret = p.data[p.pos]
	p.pos++
	return
}

// ReadBytes 读字节组
func (p *Packet) ReadBytes() (ret []byte, err error) {
	if p.pos+2 > len(p.data) {
		err = errors.New("read bytes header failed")
		return
	}
	size, _ := p.ReadU16()
	if p.pos+int(size) > len(p.data) {
		err = errors.New("read bytes data failed")
		return
	}

	ret = p.data[p.pos : p.pos+int(size)]
	p.pos += int(size)
	return
}

// ReadString 读字符
func (p *Packet) ReadString() (ret string, err error) {
	if p.pos+2 > len(p.data) {
		err = errors.New("read string header failed")
		return
	}

	size, _ := p.ReadU16()
	if p.pos+int(size) > len(p.data) {
		err = errors.New("read string data failed")
		return
	}

	bytes := p.data[p.pos : p.pos+int(size)]
	p.pos += int(size)
	ret = string(bytes)
	return
}

// ReadS8 读int8
func (p *Packet) ReadS8() (ret int8, err error) {
	_ret, _err := p.ReadByte()
	ret = int8(_ret)
	err = _err
	return
}

// ReadU16 ReadU16
func (p *Packet) ReadU16() (ret uint16, err error) {
	if p.pos+2 > len(p.data) {
		err = errors.New("read uint16 failed")
		return
	}

	buf := p.data[p.pos : p.pos+2]
	ret = uint16(buf[0])<<8 | uint16(buf[1])
	p.pos += 2
	return
}

// ReadS16 ReadS16
func (p *Packet) ReadS16() (ret int16, err error) {
	_ret, _err := p.ReadU16()
	ret = int16(_ret)
	err = _err
	return
}

// ReadU24 ReadU24
func (p *Packet) ReadU24() (ret uint32, err error) {
	if p.pos+3 > len(p.data) {
		err = errors.New("read uint24 failed")
		return
	}

	buf := p.data[p.pos : p.pos+3]
	ret = uint32(buf[0])<<16 | uint32(buf[1])<<8 | uint32(buf[2])
	p.pos += 3
	return
}

// ReadS24 ReadS24
func (p *Packet) ReadS24() (ret int32, err error) {
	_ret, _err := p.ReadU24()
	ret = int32(_ret)
	err = _err
	return
}

// ReadU32 ReadU32
func (p *Packet) ReadU32() (ret uint32, err error) {
	if p.pos+4 > len(p.data) {
		err = errors.New("read uint32 failed")
		return
	}

	buf := p.data[p.pos : p.pos+4]
	ret = uint32(buf[0])<<24 | uint32(buf[1])<<16 | uint32(buf[2])<<8 | uint32(buf[3])
	p.pos += 4
	return
}

// ReadS32 ReadS32
func (p *Packet) ReadS32() (ret int32, err error) {
	_ret, _err := p.ReadU32()
	ret = int32(_ret)
	err = _err
	return
}

// ReadU64 ReadU64
func (p *Packet) ReadU64() (ret uint64, err error) {
	if p.pos+8 > len(p.data) {
		err = errors.New("read uint64 failed")
		return
	}

	ret = 0
	buf := p.data[p.pos : p.pos+8]
	for i, v := range buf {
		ret |= uint64(v) << uint((7-i)*8)
	}
	p.pos += 8
	return
}

// ReadS64 ReadS64
func (p *Packet) ReadS64() (ret int64, err error) {
	_ret, _err := p.ReadU64()
	ret = int64(_ret)
	err = _err
	return
}

// ReadFloat32 ReadFloat32
func (p *Packet) ReadFloat32() (ret float32, err error) {
	bits, _err := p.ReadU32()
	if _err != nil {
		return float32(0), _err
	}

	ret = math.Float32frombits(bits)
	if math.IsNaN(float64(ret)) || math.IsInf(float64(ret), 0) {
		return 0, nil
	}

	return ret, nil
}

// ReadFloat64 ReadFloat64
func (p *Packet) ReadFloat64() (ret float64, err error) {
	bits, _err := p.ReadU64()
	if _err != nil {
		return float64(0), _err
	}

	ret = math.Float64frombits(bits)
	if math.IsNaN(ret) || math.IsInf(ret, 0) {
		return 0, nil
	}

	return ret, nil
}

// WriteZeros WriteZeros
func (p *Packet) WriteZeros(n int) {
	for i := 0; i < n; i++ {
		p.data = append(p.data, byte(0))
	}
}

// WriteBool WriteBool
func (p *Packet) WriteBool(v bool) {
	if v {
		p.data = append(p.data, byte(1))
	} else {
		p.data = append(p.data, byte(0))
	}
}

// WriteByte WriteByte
func (p *Packet) WriteByte(v byte) error {
	p.data = append(p.data, v)
	return nil
}

// WriteBytes WriteBytes
func (p *Packet) WriteBytes(v []byte) {
	p.WriteU16(uint16(len(v)))
	p.data = append(p.data, v...)
}

// WriteRawBytes WriteRawBytes
func (p *Packet) WriteRawBytes(v []byte) {
	p.data = append(p.data, v...)
}

// WriteString WriteString
func (p *Packet) WriteString(v string) {
	bytes := []byte(v)
	p.WriteU16(uint16(len(bytes)))
	p.data = append(p.data, bytes...)
}

// WriteS8 WriteS8
func (p *Packet) WriteS8(v int8) error {
	return p.WriteByte(byte(v))
}

// WriteU16 WriteU16
func (p *Packet) WriteU16(v uint16) {
	p.data = append(p.data, byte(v>>8), byte(v))
}

// WriteS16 WriteS16
func (p *Packet) WriteS16(v int16) {
	p.WriteU16(uint16(v))
}

// WriteU24 WriteU24
func (p *Packet) WriteU24(v uint32) {
	p.data = append(p.data, byte(v>>16), byte(v>>8), byte(v))
}

// WriteU32 WriteU32
func (p *Packet) WriteU32(v uint32) {
	p.data = append(p.data, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

// WriteS32 WriteS32
func (p *Packet) WriteS32(v int32) {
	p.WriteU32(uint32(v))
}

// WriteU64 WriteU64
func (p *Packet) WriteU64(v uint64) {
	p.data = append(p.data, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

// WriteS64 WriteS64
func (p *Packet) WriteS64(v int64) {
	p.WriteU64(uint64(v))
}

// WriteFloat32 WriteFloat32
func (p *Packet) WriteFloat32(f float32) {
	v := math.Float32bits(f)
	p.WriteU32(v)
}

// WriteFloat64 WriteFloat64
func (p *Packet) WriteFloat64(f float64) {
	v := math.Float64bits(f)
	p.WriteU64(v)
}

// Reader Reader
func Reader(data []byte) *Packet {
	return &Packet{data: data}
}

// Writer Writer
func Writer() *Packet {
	return &Packet{data: make([]byte, 0, 512)}
}
