package pack

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
)

var ErrNotEnoughBytes = errors.New("not enough bytes")

const CRC32_STEP = 4096

var crc32_buffer []byte = make([]byte, CRC32_STEP)

func CRC32IEEE(reader io.ReadSeeker, file_size int64, value *uint32) error {
	checksum := crc32.New(crc32.IEEETable)

	end_position := int((file_size/CRC32_STEP)-1) * CRC32_STEP

	for i := 0; i <= end_position; i += CRC32_STEP {
		_, err := reader.Read(crc32_buffer)
		if err != nil {
			return err
		}

		checksum.Write(crc32_buffer)
	}

	remainder := file_size % CRC32_STEP

	if remainder != 0 {
		_, err := reader.Read(crc32_buffer[:remainder])
		if err != nil {
			return err
		}

		checksum.Write(crc32_buffer[:remainder])
	}

	*value = checksum.Sum32()
	return nil
}

func ReadUInt8(reader io.ReadSeeker, value *uint8) error {
	return binary.Read(reader, binary.NativeEndian, value)
}

func WriteUInt8(writer io.WriteSeeker, value uint8) error {
	return binary.Write(writer, binary.NativeEndian, &value)
}

func ReadUInt16(reader io.ReadSeeker, isBigEndian bool, value *uint16) error {
	if isBigEndian {
		return binary.Read(reader, binary.BigEndian, value)
	} else {
		return binary.Read(reader, binary.LittleEndian, value)
	}
}

func WriteUInt16(writer io.WriteSeeker, isBigEndian bool, value uint16) error {
	if isBigEndian {
		return binary.Write(writer, binary.BigEndian, &value)
	} else {
		return binary.Write(writer, binary.LittleEndian, &value)
	}
}

func ReadUInt32(reader io.ReadSeeker, isBigEndian bool, value *uint32) error {
	if isBigEndian {
		return binary.Read(reader, binary.BigEndian, value)
	} else {
		return binary.Read(reader, binary.LittleEndian, value)
	}
}

func WriteUInt32(writer io.WriteSeeker, isBigEndian bool, value uint32) error {
	if isBigEndian {
		return binary.Write(writer, binary.BigEndian, &value)
	} else {
		return binary.Write(writer, binary.LittleEndian, &value)
	}
}

func ReadUInt64(reader io.ReadSeeker, isBigEndian bool, value *uint64) error {
	if isBigEndian {
		return binary.Read(reader, binary.BigEndian, value)
	} else {
		return binary.Read(reader, binary.LittleEndian, value)
	}
}

func WriteUInt64(writer io.WriteSeeker, isBigEndian bool, value uint64) error {
	if isBigEndian {
		return binary.Write(writer, binary.BigEndian, &value)
	} else {
		return binary.Write(writer, binary.LittleEndian, &value)
	}
}

func ReadInt8(reader io.ReadSeeker, value *int8) error {
	return binary.Read(reader, binary.NativeEndian, value)
}

func WriteInt8(writer io.WriteSeeker, value int8) error {
	return binary.Write(writer, binary.NativeEndian, &value)
}

func ReadInt16(reader io.ReadSeeker, isBigEndian bool, value *int16) error {
	if isBigEndian {
		return binary.Read(reader, binary.BigEndian, value)
	} else {
		return binary.Read(reader, binary.LittleEndian, value)
	}
}

func WriteInt16(writer io.WriteSeeker, isBigEndian bool, value int16) error {
	if isBigEndian {
		return binary.Write(writer, binary.BigEndian, &value)
	} else {
		return binary.Write(writer, binary.LittleEndian, &value)
	}
}

func ReadInt32(reader io.ReadSeeker, isBigEndian bool, value *int32) error {
	if isBigEndian {
		return binary.Read(reader, binary.BigEndian, value)
	} else {
		return binary.Read(reader, binary.LittleEndian, value)
	}
}

func WriteInt32(writer io.WriteSeeker, isBigEndian bool, value int32) error {
	if isBigEndian {
		return binary.Write(writer, binary.BigEndian, &value)
	} else {
		return binary.Write(writer, binary.LittleEndian, &value)
	}
}

func ReadInt64(reader io.ReadSeeker, isBigEndian bool, value *int64) error {
	if isBigEndian {
		return binary.Read(reader, binary.BigEndian, value)
	} else {
		return binary.Read(reader, binary.LittleEndian, value)
	}
}

func WriteInt64(writer io.WriteSeeker, isBigEndian bool, value int64) error {
	if isBigEndian {
		return binary.Write(writer, binary.BigEndian, &value)
	} else {
		return binary.Write(writer, binary.LittleEndian, &value)
	}
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

	var char uint8

	for {
		err = binary.Read(reader, binary.NativeEndian, &char)
		if err != nil {
			if err == io.EOF {
				return nil, 0, io.ErrUnexpectedEOF
			}

			return nil, 0, err
		}

		if char == 0 {
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

	return WriteUInt8(writer, 0)
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
