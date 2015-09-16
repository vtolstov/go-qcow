package qcow

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Reader struct {
	Header QCowHeader
	r      io.ReadSeeker
}

func NewReader(r io.ReadSeeker) (*Reader, error) {
	var err error
	var magick uint32
	var version uint32

	err = binary.Read(r, binary.BigEndian, &magick)
	if err != nil {
		return nil, err
	}
	if magick != 0x514649fb {
		return nil, fmt.Errorf("qcow bad magick %x != 0x514649fb", magick)
	}

	err = binary.Read(r, binary.BigEndian, &version)
	if err != nil {
		return nil, err
	}

	_, err = r.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("qcow unable to seek")
	}

	switch version {
	case 0x3:
		return NewQCow3Reader(r)
	default:
		return nil, fmt.Errorf("qcow bad version %d", version)
	}

	return nil, fmt.Errorf("qcow something wrong")
}

func NewQCow3Reader(r io.ReadSeeker) (*Reader, error) {
	var hdr QCow3Header
	var err error

	err = binary.Read(r, binary.BigEndian, &hdr)
	if err != nil {
		return nil, err
	}
	if hdr.HeaderLength < 104 {
		return nil, fmt.Errorf("qcow header too small")
	}
	return &Reader{Header: hdr, r: r}, nil
}

func (r *Reader) ReadAt(p []byte, off int64) (n int, err error) {
	hdr := r.Header.(QCow3Header)
	var fileOffset int64

	switch hdr.Version {
	case 0x3:

	}

	return r.r.ReadAt(p, fileOffset)
}
