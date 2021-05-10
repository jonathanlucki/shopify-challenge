package adapters

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// Limited testing due to personal life time constraints
func TestGetImageUrl(t *testing.T) {
	const storageURL string = "https://test.com/"
	const id string = "abcd"
	const IdURL string = "https://test.com/abcd.png"

	// test when $STORAGE_URL not set - should return error
	result, err := GetImageUrl(id)
	require.Equal(t, result, "")
	require.NotNil(t, err)

	// test proper result - should not return error
	os.Setenv("STORAGE_URL", storageURL)
	result, err = GetImageUrl(id)
	require.Equal(t, result, IdURL)
	require.Nil(t, err)
}
