package core

import (
	"bytes"
	"testing"
)

func FuzzBitConversion(f *testing.F) {
	f.Add(byte(8))
	f.Add(byte(0))
	f.Add(byte(128))
	f.Add(byte(255))

	binary := Binary{}
	f.Fuzz(func(t *testing.T, a byte) {
		bitArray := binary.ByteToBitArray(a)
		if b, err := binary.BitArrayToByte(bitArray); err != nil {
			t.Error(err)
		} else if b != a {
			t.Errorf("expected %v to equal %v", b, a)
		}
	})
}

func FuzzBitPadding(f *testing.F) {
	f.Add([]byte("Hello world"), uint64(32))

	tinyBlock, err := RandomBytes(16)
	if err != nil {
		f.Error(err)
	}
	f.Add(tinyBlock, uint64(64))

	smallBlock, err := RandomBytes(48)
	if err != nil {
		f.Error(err)
	}
	f.Add(smallBlock, uint64(64))

	midBlock, err := RandomBytes(256)
	if err != nil {
		f.Error(err)
	}
	f.Add(midBlock, uint64(32))

	largeBlock, err := RandomBytes(512)
	if err != nil {
		f.Error(err)
	}
	f.Add(largeBlock, uint64(64))

	superLargeBlock, err := RandomBytes(2560)
	if err != nil {
		f.Error(err)
	}
	f.Add(superLargeBlock, uint64(32))

	binary := Binary{}
	f.Fuzz(func(t *testing.T, a []byte, b uint64) {
		paddedBytes, err := binary.PadBytes(a, b)
		if err != nil {
			t.Error(err)
		}

		unpaddedBytes, err := binary.UnpadBytes(paddedBytes)
		if err != nil {
			t.Error(err)
		} else if !bytes.Equal(unpaddedBytes, a) {
			t.Errorf("expected %v to match %v", unpaddedBytes, a)
		}
	})
}

func TestExternalByteUnpadding(t *testing.T) {
	binary := Binary{}
	data := []byte{5, 1<<1}
	dataClone := []byte{5, 1<<1}

	lastByte := binary.ByteToBitArray(dataClone[1])
	lastByte[7] = 1
	compressedByte, err := binary.BitArrayToByte(lastByte)
	if err != nil {
		t.Error(err)
	}
	dataClone[1] = compressedByte

	unpaddedBytes, err := binary.UnpadBytes(dataClone)
	if err != nil {
		t.Error(err)
	}

	if unpaddedBytes[1] != data[1] {
		t.Errorf("expected %v to equal %v", unpaddedBytes[1], data[1])
	}
}
