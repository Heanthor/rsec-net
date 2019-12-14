package maputils

import (
	"testing"

	cmap "github.com/orcaman/concurrent-map"
	"github.com/stretchr/testify/assert"
)

func Test_ComputeHash(t *testing.T) {
	type testStruct struct {
		field1 string
		field2 []int
	}
	cmap1 := cmap.New()
	cmap1.Set("key1", testStruct{"a", []int{1, 2, 3}})
	hash1, _ := ComputeHash(cmap1)
	cmap2 := cmap.New()
	cmap2.Set("key1", testStruct{"a", []int{1, 2, 3}})
	hash2, _ := ComputeHash(cmap2)

	assert.Equal(t, hash1, hash2)

	cmap2.Set("key1", testStruct{"a", []int{1, 3, 3}})
	hash2, _ = ComputeHash(cmap2)

	assert.NotEqual(t, hash1, hash2)
}
