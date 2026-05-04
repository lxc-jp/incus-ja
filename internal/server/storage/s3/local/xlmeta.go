package local

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// xlMetaMagic is the four-byte prefix of every xl.meta file written by minio.
var xlMetaMagic = [4]byte{'X', 'L', '2', ' '}

// xlInlineData parses a minio xl.meta file and returns the inline data blob
// for the current object version, if the object is stored inline.
//
// minio's xl-storage-format wraps the inline data inside a msgpack
// structure whose exact shape varies between format revisions. Rather than
// decode the entire metadata structure, this function scans the file for
// minio's stable inline-section signature:
//
//	byte 0:  0x01  (inline-section format version)
//	byte 1+: msgpack map[string][]byte (version-id → object bytes)
//
// The first byte position whose tail parses as a non-empty msgpack map
// with a str key and a bin value is taken as the inline data section,
// and the value of the first map entry is returned.
//
// Returns (nil, nil) if the file is well-formed but contains no inline
// data section (e.g. the object's data is stored externally as part files).
func xlInlineData(b []byte) ([]byte, error) {
	if len(b) < 12 {
		return nil, errors.New("xl.meta truncated")
	}

	if [4]byte{b[0], b[1], b[2], b[3]} != xlMetaMagic {
		return nil, errors.New("xl.meta bad magic")
	}

	// Search after the 8-byte header for the inline section's signature.
	// The CRC trailer occupies the last 4 bytes; nothing useful starts
	// there.
	for p := 8; p < len(b)-4; p++ {
		if b[p] != 0x01 {
			continue
		}

		data, ok := tryReadInlineSection(b[p+1:])
		if ok {
			return data, nil
		}
	}

	return nil, nil
}

// tryReadInlineSection attempts to read a msgpack map[string][]byte at the
// start of b and returns the value of the first entry. It returns (nil,
// false) on any decoding error.
func tryReadInlineSection(b []byte) ([]byte, bool) {
	r := newMsgpReader(b)

	count, err := r.readMapLen()
	if err != nil || count == 0 {
		return nil, false
	}

	_, err = r.readStr()
	if err != nil {
		return nil, false
	}

	data, err := r.readBin()
	if err != nil {
		return nil, false
	}

	return data, true
}

// msgpReader is a minimal msgpack reader supporting only the type families
// used by minio xl.meta: bin, str, fixarray, fixmap, array16/32, map16/32.
type msgpReader struct {
	b []byte
	i int
}

func newMsgpReader(b []byte) *msgpReader { return &msgpReader{b: b} }

func (r *msgpReader) need(n int) error {
	if r.i+n > len(r.b) {
		return errors.New("msgpack short read")
	}

	return nil
}

func (r *msgpReader) readByte() (byte, error) {
	err := r.need(1)
	if err != nil {
		return 0, err
	}

	c := r.b[r.i]
	r.i++
	return c, nil
}

func (r *msgpReader) readUint(n int) (uint32, error) {
	err := r.need(n)
	if err != nil {
		return 0, err
	}

	var v uint32
	switch n {
	case 1:
		v = uint32(r.b[r.i])
	case 2:
		v = uint32(binary.BigEndian.Uint16(r.b[r.i:]))
	case 4:
		v = binary.BigEndian.Uint32(r.b[r.i:])
	}

	r.i += n
	return v, nil
}

func (r *msgpReader) readMapLen() (int, error) {
	c, err := r.readByte()
	if err != nil {
		return 0, err
	}

	switch {
	case c >= 0x80 && c <= 0x8f:
		return int(c & 0x0f), nil
	case c == 0xde:
		v, err := r.readUint(2)
		return int(v), err
	case c == 0xdf:
		v, err := r.readUint(4)
		return int(v), err
	}

	return 0, fmt.Errorf("not a map (0x%02x)", c)
}

func (r *msgpReader) readBin() ([]byte, error) {
	c, err := r.readByte()
	if err != nil {
		return nil, err
	}

	var n uint32

	switch c {
	case 0xc4:
		n, err = r.readUint(1)
	case 0xc5:
		n, err = r.readUint(2)
	case 0xc6:
		n, err = r.readUint(4)
	default:
		return nil, fmt.Errorf("not a bin (0x%02x)", c)
	}

	if err != nil {
		return nil, err
	}

	err = r.need(int(n))
	if err != nil {
		return nil, err
	}

	data := r.b[r.i : r.i+int(n)]
	r.i += int(n)
	return data, nil
}

func (r *msgpReader) readStr() (string, error) {
	c, err := r.readByte()
	if err != nil {
		return "", err
	}

	var n uint32

	switch {
	case c >= 0xa0 && c <= 0xbf:
		n = uint32(c & 0x1f)
	case c == 0xd9:
		n, err = r.readUint(1)
	case c == 0xda:
		n, err = r.readUint(2)
	case c == 0xdb:
		n, err = r.readUint(4)
	default:
		return "", fmt.Errorf("not a str (0x%02x)", c)
	}

	if err != nil {
		return "", err
	}

	err = r.need(int(n))
	if err != nil {
		return "", err
	}

	s := string(r.b[r.i : r.i+int(n)])
	r.i += int(n)
	return s, nil
}
