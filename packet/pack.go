package packet

import (
	"log"
	"reflect"
)

// FastPack FastPack
type FastPack interface {
	Pack(w *Packet)
}

// Pack export struct fields with packet writer.
func Pack(tos int16, tbl interface{}, writer *Packet) []byte {
	// create writer if not specified
	if writer == nil {
		writer = Writer()
	}

	// write protocol number
	writer.WriteS16(tos)

	// is the table nil?
	if tbl == nil {
		return writer.Data()
	}

	// fastpack
	if fastpack, ok := tbl.(FastPack); ok {
		fastpack.Pack(writer)
		return writer.Data()
	}

	// pack by reflection
	err := _pack(reflect.ValueOf(tbl), writer)
	if err != nil {
		return nil
	}

	// return byte array
	return writer.Data()
}

// _pack export struct fields with packet writer.
func _pack(v reflect.Value, writer *Packet) error {
	switch v.Kind() {
	case reflect.Bool:
		writer.WriteBool(v.Bool())
	case reflect.Uint8:
		return writer.WriteByte(byte(v.Uint()))
	case reflect.Uint16:
		writer.WriteU16(uint16(v.Uint()))
	case reflect.Uint32:
		writer.WriteU32(uint32(v.Uint()))
	case reflect.Uint64:
		writer.WriteU64(v.Uint())
	case reflect.Int16:
		writer.WriteS16(int16(v.Int()))
	case reflect.Int32:
		writer.WriteS32(int32(v.Int()))
	case reflect.Int64:
		writer.WriteS64(v.Int())
	case reflect.Float32:
		writer.WriteFloat32(float32(v.Float()))
	case reflect.Float64:
		writer.WriteFloat64(float64(v.Float()))
	case reflect.String:
		writer.WriteString(v.String())
	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			return nil
		}
		return _pack(v.Elem(), writer)
	case reflect.Slice:
		if bs, ok := v.Interface().([]byte); ok { // special treat for []bytes
			writer.WriteBytes(bs)
		} else {
			l := v.Len()
			writer.WriteU16(uint16(l))
			for i := 0; i < l; i++ {
				if err := _pack(v.Index(i), writer); err != nil {
					return err
				}
			}
		}
	case reflect.Struct:
		numFields := v.NumField()
		for i := 0; i < numFields; i++ {
			if err := _pack(v.Field(i), writer); err != nil {
				return err
			}
		}
	default:
		log.Println("cannot pack type:", v)
	}
	return nil
}
