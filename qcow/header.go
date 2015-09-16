package qcow

// https://github.com/qemu/QEMU/blob/master/docs/specs/qcow2.txt
// https://github.com/libyal/libqcow/blob/master/documentation/QEMU%20Copy-On-Write%20file%20format.asciidoc

type QCowHeader interface {
}

type QCow1Header struct {
	Magick  uint32
	Version uint32

	BackingFileOffset uint64
	BackingFileSize   uint32

	Mtime uint32

	Size uint64

	ClusterBits uint8
	L2Bits      uint8

	Unused1 uint16

	CryptMethod uint32

	L1TableOffset uint64
}

func (h *QCow1Header) clusterBlockSize() uint64 {
	return (1 << h.ClusterBits)
}

func (h *QCow1Header) l2TableSize() uint64 {
	return (1 << h.L2Bits) * 8
}

func (h *QCow1Header) l1TableSize() uint64 {
	l1TableSize := h.clusterBlockSize() * (1 << h.L2Bits)

	if h.Size%l1TableSize != 0 {
		l1TableSize = (h.Size / l1TableSize) + 1
	} else {
		l1TableSize = h.Size / l1TableSize
	}
	return l1TableSize * 8
}

func (h *QCow2Header) clusterBlockSize() uint64 {
	return uint64(1 << h.ClusterBits)
}

func (h *QCow2Header) l2TableBits() uint64 {
	return uint64(h.ClusterBits - 3)
}

func (h *QCow2Header) l2TableSize() uint64 {
	return (1 << h.l2TableBits()) * 8
}

func (h *QCow2Header) l1TableSize() uint64 {
	return uint64(h.L1Size * 8)
}

func (h *QCow2Header) l1TableIndexBitShift() uint64 {
	return uint64(h.ClusterBits) + h.l2TableBits()
}

func (h *QCow2Header) l1TableIndex(offset int64) uint64 {
	return (uint64(offset) & uint64(0x3fffffffffffffff)) >> h.l1TableIndexBitShift()
}

func (h *QCow2Header) l2TableIndexBitMask() uint64 {
	return ^(uint64(0xffffffffffffffff) << h.l2TableSize())
}

func (h *QCow2Header) l2TableIndex(offset int64) uint64 {
	return (uint64(offset) >> h.ClusterBits) >> h.l2TableIndexBitMask()
}

func (h *QCow2Header) compSizeBitShift() uint64 {
	return uint64(62 - (h.ClusterBits - 8))
}

func (h *QCow2Header) clusterBlockOffset() uint64 {
	return ^(uint64(0xffffffffffffffff) << h.compSizeBitShift())
}

func (h *QCow2Header) compBlockSize() uint64 {
	return (((h.clusterBlockOffset() & 0x3fffffffffffffff) >> h.compSizeBitShift()) + 1) * 512
}

type QCow2Header struct {
	Magick  uint32
	Version uint32

	BackingFileOffset uint64
	BackingFileSize   uint32

	ClusterBits           uint32
	Size                  uint64
	CryptMethod           uint32
	L1Size                uint32
	L1TableOffset         uint64
	RefcountTableOffset   uint64
	RefcountTableClusters uint32
	NbSnapshots           uint32
	SnapshotsOffset       uint64
}

type QCow3Header struct {
	Magick  uint32
	Version uint32

	BackingFileOffset uint64
	BackingFileSize   uint32

	ClusterBits           uint32
	Size                  uint64
	CryptMethod           uint32
	L1Size                uint32
	L1TableOffset         uint64
	RefcountTableOffset   uint64
	RefcountTableClusters uint32
	NbSnapshots           uint32
	SnapshotsOffset       uint64

	IncompatibleFeatures uint64
	CompatibleFeatures   uint64
	AutoclearFeatures    uint64
	RefcountOrder        uint32
	HeaderLength         uint32
}
