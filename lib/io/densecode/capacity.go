package densecode

func calculateMaxCapacity(errorLevel, bitsPerModule int) int {
	if bitsPerModule <= 0 {
		bitsPerModule = 3
	}

	maxModules := 1000*1000 - 3*7*7 - 2*(1000-14)

	bitsAvailable := maxModules * bitsPerModule

	bytesAvailable := bitsAvailable / 8

	redundancy := (errorLevel + 1) * 8
	netCapacity := (bytesAvailable*100)/(100+redundancy*100/bytesAvailable) - 20

	return netCapacity / 2
}

func calculateOptimalSegmentSize(errorLevel, bitsPerModule int) int {
	baseSize := 32 * 1024

	switch errorLevel {
	case 0:
		return baseSize
	case 1:
		return baseSize
	case 2:
		return baseSize * 3 / 4
	case 3:
		return baseSize / 2
	default:
		return baseSize
	}
}
