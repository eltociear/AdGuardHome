package aghalg

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSortedMap(t *testing.T) {
	var m SortedMap[string, int]

	letters := []string{}
	for i := 0; i < 10; i++ {
		r := string('a' + rune(i))
		letters = append(letters, r)
	}

	t.Run("create_and_fill", func(t *testing.T) {
		m = NewSortedMap[string, int](strings.Compare)

		nums := []int{}
		for i, r := range letters {
			m.Set(r, i)
			nums = append(nums, i)
		}

		gotLetters := []string{}
		gotNums := []int{}
		m.Range(func(k string, v int) bool {
			gotLetters = append(gotLetters, k)
			gotNums = append(gotNums, v)

			return true
		})

		assert.Equal(t, letters, gotLetters)
		assert.Equal(t, nums, gotNums)
		assert.Equal(t, nums[0], m.Get(letters[0]))
	})

	t.Run("clear", func(t *testing.T) {
		lastLetter := letters[len(letters)-1]
		m.Del(lastLetter)
		assert.Equal(t, 0, m.Get(lastLetter))

		m.Clear()

		gotLetters := []string{}
		m.Range(func(k string, _ int) bool {
			gotLetters = append(gotLetters, k)

			return true
		})

		assert.Len(t, gotLetters, 0)
	})
}
