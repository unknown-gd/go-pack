package pack

import (
	"encoding/binary"
	"errors"
	"io"
)

var ErrNotEnoughBytes = errors.New("not enough bytes")

var long []byte = make([]byte, 8)
var integer []byte = make([]byte, 4)
var short []byte = make([]byte, 2)
var char []byte = make([]byte, 1)

func ReadUInt8(reader io.ReadSeeker) (uint8, error) {
	_, err := reader.Read(char[:1])
	if err != nil {
		return 0, err
	}

	return char[0], nil
}

func WriteUInt8(writer io.WriteSeeker, value uint8) error {
	char[0] = value
	_, err := writer.Write(char[:1])
	return err
}

func ReadUInt16(reader io.ReadSeeker, isBigEndian bool) (uint16, error) {
	_, err := reader.Read(short[:2])
	if err != nil {
		return 0, err
	}

	if isBigEndian {
		return binary.BigEndian.Uint16(short[:2]), nil
	} else {
		return binary.LittleEndian.Uint16(short[:2]), nil
	}
}

func WriteUInt16(writer io.WriteSeeker, value uint16, isBigEndian bool) error {
	if isBigEndian {
		binary.BigEndian.PutUint16(short[:2], value)
	} else {
		binary.LittleEndian.PutUint16(short[:2], value)
	}

	_, err := writer.Write(short[:2])
	return err
}

func ReadUInt32(reader io.ReadSeeker, isBigEndian bool) (uint32, error) {
	_, err := reader.Read(integer[:4])
	if err != nil {
		return 0, err
	}

	if isBigEndian {
		return binary.BigEndian.Uint32(integer[:4]), nil
	} else {
		return binary.LittleEndian.Uint32(integer[:4]), nil
	}
}

func WriteUInt32(writer io.WriteSeeker, value uint32, isBigEndian bool) error {
	if isBigEndian {
		binary.BigEndian.PutUint32(integer[:4], value)
	} else {
		binary.LittleEndian.PutUint32(integer[:4], value)
	}

	_, err := writer.Write(integer[:4])
	return err
}

func ReadUInt64(reader io.ReadSeeker, isBigEndian bool) (uint64, error) {
	_, err := reader.Read(long[:8])
	if err != nil {
		return 0, err
	}

	if isBigEndian {
		return binary.BigEndian.Uint64(long[:8]), nil
	} else {
		return binary.LittleEndian.Uint64(long[:8]), nil
	}
}

func WriteUInt64(writer io.WriteSeeker, value uint64, isBigEndian bool) error {
	if isBigEndian {
		binary.BigEndian.PutUint64(long[:8], value)
	} else {
		binary.LittleEndian.PutUint64(long[:8], value)
	}

	_, err := writer.Write(long[:8])
	return err
}

func ReadInt8(reader io.ReadSeeker) (int8, error) {
	x, err := ReadUInt8(reader)
	return int8(x), err
}

func WriteInt8(writer io.WriteSeeker, value int8) error {
	return WriteUInt8(writer, uint8(value))
}

func ReadInt16(reader io.ReadSeeker, isBigEndian bool) (int16, error) {
	x, err := ReadUInt16(reader, isBigEndian)
	return int16(x), err
}

func WriteInt16(writer io.WriteSeeker, value int16, isBigEndian bool) error {
	return WriteUInt16(writer, uint16(value), isBigEndian)
}

func ReadInt32(reader io.ReadSeeker, isBigEndian bool) (int32, error) {
	x, err := ReadUInt32(reader, isBigEndian)
	return int32(x), err
}

func WriteInt32(writer io.WriteSeeker, value int32, isBigEndian bool) error {
	return WriteUInt32(writer, uint32(value), isBigEndian)
}

func ReadInt64(reader io.ReadSeeker, isBigEndian bool) (int64, error) {
	x, err := ReadUInt64(reader, isBigEndian)
	return int64(x), err
}

func WriteInt64(writer io.WriteSeeker, value int64, isBigEndian bool) error {
	return WriteUInt64(writer, uint64(value), isBigEndian)
}

func ReadFixedBytes(reader io.ReadSeeker, length *uint32) ([]byte, uint32, error) {
	var read_length uint32

	if length == nil {
		read_length = 255
	} else {
		read_length = *length
	}

	data := make([]byte, read_length)

	_, err := reader.Read(data[:read_length])
	if err != nil {
		return nil, 0, err
	}

	return data, read_length, nil
}

func WriteFixedBytes(writer io.WriteSeeker, data []byte, length *uint32) error {
	var write_length uint32

	if length == nil {
		write_length = uint32(len(data))
	} else {
		write_length = *length
	}

	_, err := writer.Write(data[:write_length])
	if err != nil {
		return err
	}

	return nil
}

func ReadCountedBytes(reader io.ReadSeeker, count *uint32, isBigEndian bool) ([]byte, uint32, error) {
	var byte_count uint32

	if count == nil {
		byte_count = 2 // short (2 bytes/uint16 ~= 65535 bytes)
	} else {
		byte_count = *count
	}

	byte_buffer := make([]byte, byte_count)

	_, err := reader.Read(byte_buffer)
	if err != nil {
		return nil, 0, err
	}

	var read_length uint32 = 0

	if isBigEndian {
		read_length = binary.BigEndian.Uint32(byte_buffer)
	} else {
		read_length = binary.LittleEndian.Uint32(byte_buffer)
	}

	return ReadFixedBytes(reader, &read_length)
}

func WriteCountedBytes(writer io.WriteSeeker, data []byte, count *uint32, length *uint32, isBigEndian bool) error {
	buffer_size := uint32(len(data))

	var byte_count uint32

	if count == nil {
		byte_count = 2 // short (2 bytes/uint16 ~= 65535 bytes)
	} else {
		byte_count = *count
	}

	var write_length uint32

	if length == nil {
		write_length = buffer_size
	} else {
		write_length = *length

		if write_length > buffer_size {
			return ErrNotEnoughBytes
		}
	}

	if write_length > (2 ^ byte_count - 1) {
		return ErrNotEnoughBytes
	}

	byte_buffer := make([]byte, byte_count)

	if isBigEndian {
		binary.BigEndian.PutUint32(byte_buffer, write_length)
	} else {
		binary.LittleEndian.PutUint32(byte_buffer, write_length)
	}

	_, err := writer.Write(byte_buffer)
	if err != nil {
		return err
	}

	_, err = writer.Write(data[:write_length])
	if err != nil {
		return err
	}

	return nil
}

func ReadNullTerminatedBytes(reader io.ReadSeeker) ([]byte, uint32, error) {
	start_position, err := reader.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, 0, err
	}

	var end_position int64
	var read_length uint32

	for {
		_, err := reader.Read(char[:1])
		if err != nil {
			if err == io.EOF {
				return nil, 0, io.ErrUnexpectedEOF
			}

			return nil, 0, err
		}

		if char[0] == 0 {
			end_position, err = reader.Seek(0, io.SeekCurrent)
			if err != nil {
				return nil, 0, err
			}

			end_position = end_position - 1

			if start_position == end_position {
				read_length = 0
			} else {
				read_length = uint32(end_position - start_position)
			}

			data := make([]byte, read_length)

			if read_length == 0 {
				return data, uint32(read_length), nil
			}

			_, err = reader.Seek(start_position, io.SeekStart)
			if err != nil {
				return nil, 0, err
			}

			_, err = reader.Read(data)
			if err != nil {
				return nil, 0, err
			}

			reader.Seek(1, io.SeekCurrent)
			return data, uint32(read_length), nil
		}
	}
}

func WriteNullTerminatedBytes(writer io.WriteSeeker, data []byte, length *uint32) error {
	var write_length uint32

	if length == nil {
		write_length = uint32(len(data))
	} else {
		write_length = *length
	}

	_, err := writer.Write(data[:write_length])
	if err != nil {
		return err
	}

	char[0] = 0

	_, err = writer.Write(char[:1])
	if err != nil {
		return err
	}

	return nil
}

func ReadFixedString(reader io.ReadSeeker, length *uint32) (string, uint32, error) {
	data, read_length, err := ReadFixedBytes(reader, length)
	if err != nil {
		return "", 0, err
	}

	return string(data), read_length, nil
}

func WriteFixedString(writer io.WriteSeeker, str string, length *uint32) error {
	return WriteFixedBytes(writer, []byte(str), length)
}

func ReadCountedString(reader io.ReadSeeker, count *uint32, isBigEndian bool) (string, uint32, error) {
	data, read_length, err := ReadCountedBytes(reader, count, isBigEndian)
	if err != nil {
		return "", 0, err
	}

	return string(data), read_length, nil
}

func WriteCountedString(writer io.WriteSeeker, str string, count *uint32, length *uint32, isBigEndian bool) error {
	return WriteCountedBytes(writer, []byte(str), count, length, isBigEndian)
}

func ReadNullTerminatedString(reader io.ReadSeeker) (string, uint32, error) {
	data, read_length, err := ReadNullTerminatedBytes(reader)
	if err != nil {
		return "", 0, err
	}

	return string(data), read_length, nil
}

func WriteNullTerminatedString(writer io.WriteSeeker, str string) error {
	return WriteNullTerminatedBytes(writer, []byte(str), nil)
}
