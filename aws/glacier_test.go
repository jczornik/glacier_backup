package aws

import (
	"testing"
)

func TestPartCannotBeSmallerThanMinPartSize(t *testing.T) {
	// Given
	size := int64(minPartSize - 1)

	// When
	_, err := computePartSize(size)

	// Then
	if err == nil {
		t.Error("Should return an error when archive is too small")
	}
}

func TestCannotExceedMaxParts(t *testing.T) {
	// Given
	size := int64(maxParts*maxPartSize + 1)

	// When
	_, err := computePartSize(size)

	// Then
	if err == nil {
		t.Error("Should return an error when archive is too big")
	}
}

func TestUseOneMBParts(t *testing.T) {
	// Given
	size := int64(minPartSize * 2)

	// When
	partSize, err := computePartSize(size)

	// Then
	if err != nil {
		t.Error("Should return part size without an error")
	}

	if partSize != minPartSize {
		t.Errorf("Expected min part size %d but received %d", minPartSize, partSize)
	}
}

func TestCorrectPartSizeWhenLastSmaller(t *testing.T) {
	// Given
	size := int64(minPartSize*2 + 1)

	// When
	partSize, err := computePartSize(size)

	// Then
	if err != nil {
		t.Error("Should return part size without an error")
	}

	if partSize != minPartSize {
		t.Errorf("Expected min part size %d but received %d", minPartSize, partSize)
	}
}

func TestReturnMaxPartSizeWhenArchiveSizeIsMax(t *testing.T) {
	// Given
	size := int64(maxParts * maxPartSize)

	// When
	partSize, err := computePartSize(size)

	// Then
	if err != nil {
		t.Error("Should return max part size without an error")
	}

	if partSize != maxPartSize {
		t.Errorf("Should return max part size (%d) but received %d", maxPartSize, partSize)
	}
}
