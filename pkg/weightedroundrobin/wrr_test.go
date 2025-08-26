package wrrsimple

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWRRSimple(t *testing.T) {
	ast := assert.New(t)
	wrr := New[string]()
	ast.NotNil(wrr)
	wrrimpl, ok := wrr.(*wrrSimple[string])
	ast.True(ok)
	ast.NotNil(wrrimpl)
	err := wrr.Add("1", 10)
	ast.Nil(err)
	err = wrr.Add("2", 5)
	ast.Nil(err)
	err = wrr.Add("3", 2)
	ast.Nil(err)

	list := make(map[string]int)
	list["1"] = 0
	list["2"] = 0
	list["3"] = 0

	for x := range 100 {
		i, err := wrr.GetNext()
		ast.Nil(err)
		ast.NotNil(i)
		list[*i]++
		if x%1000 == 0 {
			for _, i := range wrrimpl.items {
				t.Logf("actual weights: %s: %d", i.id, i.actual)
			}
		}
	}

	t.Logf("weights: 1: %d, 2: %d, 3: %d", list["1"], list["2"], list["3"])
	ast.GreaterOrEqual(list["1"], 2*(list["2"]-1))
	ast.GreaterOrEqual(list["1"], 5*(list["3"]-1))
}

func TestWRRSimpleEquality(t *testing.T) {
	ast := assert.New(t)
	wrr := New[string]()
	ast.NotNil(wrr)
	wrrimpl, ok := wrr.(*wrrSimple[string])
	ast.True(ok)
	ast.NotNil(wrrimpl)

	list := make(map[string]int)
	for x := range 10 {
		id := fmt.Sprintf("%d", x)
		err := wrr.Add(id, 10)
		ast.Nil(err)
		list[id] = 0
	}

	for range 10000 {
		i, err := wrr.GetNext()
		ast.Nil(err)
		ast.NotNil(i)
		list[*i]++
	}

	for _, v := range list {
		ast.Equal(1000, v)
	}
}

func TestNoItems(t *testing.T) {
	ast := assert.New(t)

	wss := New[string]()
	ast.NotNil(wss)

	_, err := wss.GetNext()
	ast.ErrorIs(ErrNoItems, err)
}
