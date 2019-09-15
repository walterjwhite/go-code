package disk

import (
	//"os"
	"fmt"
	"math"
	//"strconv"
	"github.com/walterjwhite/go-application/libraries/logging"
	"syscall"
)

const decimalPlaces = 2
const roundOn = .5

var units [5]string

func init() {
	units[0] = "B"
	units[1] = "KB"
	units[2] = "MB"
	units[3] = "GB"
	units[4] = "TB"
}

func Usage(path string) *Disk {
	var stat syscall.Statfs_t

	logging.Panic(syscall.Statfs(path, &stat))

	fmt.Printf("blocks: %v\n", stat.Blocks)
	fmt.Printf("free blocks: %v\n", stat. /*Bfree*/ Bavail)

	usagePercentage := uint(100 * (1.0 * (stat.Blocks - stat. /*Bfree*/ Bavail) / stat.Blocks))

	// TODO: configure units (bytes -> GB)
	freeBytes := stat. /*Bfree*/ Bavail * uint64(stat.Bsize)
	prettyFreeBytes, prettyFreeUnits := Size(freeBytes)

	free := fmt.Sprintf("%v %v", prettyFreeBytes, prettyFreeUnits)

	//Bavail

	return &Disk{MountPoint: path, UsagePercentage: usagePercentage, Free: free}
}

func Size(sizeInBytes uint64) (float64, string) {
	base := math.Log(float64(sizeInBytes)) / math.Log(1024)
	getSize := round(math.Pow(1024, base-math.Floor(base)), roundOn, decimalPlaces)
	getSuffix := units[int(math.Floor(base))]

	//return strconv.FormatFloat(getSize, 'f', -1, 64), getSuffix
	return getSize, getSuffix
}

func round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}

	newVal = round / pow
	return
}
