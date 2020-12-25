package pre_generator

import (
	"encoding/binary"
	"github.com/byxor/NeverScript/compiler"
	"io/ioutil"
	"log"
	"strings"
)

type PreSpec []PreSpecItem
type PreSpecItem struct {
	PathOnDisk    string
	PathInsidePre string
}

func ParsePreSpec(preSpecPath string) (preSpec []PreSpecItem) {
	fileBytes, err := ioutil.ReadFile(preSpecPath)
	if err != nil {
		log.Fatal(err)
	}
	text := strings.Replace(string(fileBytes), "\r", "", -1)
	lines := strings.Split(text, "\n")

	// Parser has 2 states
	var state int
	state_readingPathOnDisk := 0
	state_readingPathInsidePre := 1

	// Parse each item in the pre spec
	var pathOnDisk string
	var pathInsidePre string
	for _, line := range lines {
		switch state {
		case state_readingPathOnDisk:
			if line == "" {
				continue
			}
			pathOnDisk = line
			state = state_readingPathInsidePre
		case state_readingPathInsidePre:
			pathInsidePre = line
			preSpec = append(preSpec, PreSpecItem{
				PathOnDisk:    pathOnDisk,
				PathInsidePre: pathInsidePre,
			})
			state = state_readingPathOnDisk
		default:
			log.Println("Reached unexpected state when parsing pre spec")
			break
		}
	}

	return preSpec
}

func MakePre(preSpec PreSpec) []byte {
	pre := make([]byte, 25000000) // FIXME(brandon): Arbitrarily-sized buffer, could crash if low on RAM or if Pre is too large.

	var globalHeader struct {
		Size          uint32
		Version       uint32
		NumberOfItems uint32
	}
	globalHeader.Version = 0xABCD0003

	offset := uint32(12)
	for _, preSpecItem := range preSpec {
		// Read next pre item into memory
		fileBytes, err := ioutil.ReadFile(preSpecItem.PathOnDisk)
		if err != nil {
			log.Fatal(err)
		}
		fileLength := uint32(len(fileBytes))

		// Calculate header values
		var preItemHeader struct {
			InflatedSize          uint32
			DeflatedSize          uint32
			PathInsidePreLength   uint32
			PathInsidePre         string
			PathInsidePreChecksum uint32
		}
		preItemHeader.InflatedSize = uint32(len(fileBytes))
		preItemHeader.DeflatedSize = 0
		preItemHeader.PathInsidePreLength = uint32(len(preSpecItem.PathInsidePre) + 1)
		preItemHeader.PathInsidePre = preSpecItem.PathInsidePre
		preItemHeader.PathInsidePreChecksum = compiler.StringToChecksum(preSpecItem.PathInsidePre)

		// Align end of pathInsidePre string with nearest 4th byte
		x := offset + 12 + preItemHeader.PathInsidePreLength
		for x%4 != 0 {
			x++
		}
		preItemHeader.PathInsidePreLength = x - offset - 12 // A little bit of algebra

		// Write header
		binary.LittleEndian.PutUint32(pre[offset:], preItemHeader.InflatedSize)
		binary.LittleEndian.PutUint32(pre[offset+4:], preItemHeader.DeflatedSize)
		binary.LittleEndian.PutUint32(pre[offset+8:], preItemHeader.PathInsidePreLength)
		binary.LittleEndian.PutUint32(pre[offset+12:], preItemHeader.PathInsidePreChecksum)
		offset += 16
		copy(pre[offset:], preItemHeader.PathInsidePre)
		offset += preItemHeader.PathInsidePreLength
		pre[offset-1] = 0

		// Align offset with nearest 4th byte
		for offset%4 != 0 {
			offset++
		}

		// Write pre item contents
		copy(pre[offset:], fileBytes)
		offset += fileLength

		// Align offset with nearest 4th byte
		for offset%4 != 0 {
			offset++
		}

		// Update global header data
		globalHeader.NumberOfItems++
	}
	globalHeader.Size = offset

	// Write global header at start of pre (now that we have all the necessary info)
	binary.LittleEndian.PutUint32(pre, globalHeader.Size)
	binary.LittleEndian.PutUint32(pre[4:], globalHeader.Version)
	binary.LittleEndian.PutUint32(pre[8:], globalHeader.NumberOfItems)

	return pre[:offset]
}
