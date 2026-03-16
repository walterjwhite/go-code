package densecode

import (
	"crypto/sha256"
)

func removeErrorCorrectionBytes(data []byte, errorLevel int) []byte {
	if len(data) <= 8 {
		return data
	}

	redundancy := (errorLevel + 1) * 8
	if len(data) < redundancy {
		return data
	}

	dataOnly := data[:len(data)-redundancy]

	for i := range redundancy {
		parity := byte(0)
		for j := i; j < len(dataOnly); j += redundancy {
			parity ^= dataOnly[j]
		}
		if parity != data[len(dataOnly)+i] {
			break
		}
	}

	return dataOnly
}

func (c *Configuration) calculateDataChecksum(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:8]
}

func (c *Configuration) checksumsMatch(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
