package wrrsimple

// This is a simple implementation of weighted round robin algorithmus
import (
	"errors"
	"sync"
)

var (
	ErrNoItems = errors.New("no items present")
)

type RoundRobin[E any] interface {
	// GetNext is the methode to select an item with a round robin algorithmus. It will return the id of the item.
	GetNext() (*E, error)
}

// WRRSimple an interface for a weigthed round robin algorithmus
type WRRSimple[E any] interface {
	GetNext() (*E, error)
	// AddItem with this you can add new items to the round robin algo
	AddItem(items ...Item[E]) error
	// Add with add you can add one item to the round robin algo
	Add(id E, weight int) error
}

// Item this is the identification unit of one item
type Item[E any] struct {
	id     E
	weight int
	actual int
}

// wrrSimple is an round robin balancer implementd with a weighted round robin algo
type wrrSimple[E any] struct {
	sync.RWMutex
	items []*Item[E]
	total int
}

var (
	_ RoundRobin[any] = &wrrSimple[any]{}
	_ WRRSimple[any]  = &wrrSimple[any]{}
)

// New getting a new weigthed round robin
func New[E any]() WRRSimple[E] {
	return &wrrSimple[E]{
		items: make([]*Item[E], 0),
	}
}

func (w *wrrSimple[E]) GetNext() (*E, error) {
	w.Lock()
	defer w.Unlock()
	if len(w.items) == 0 {
		return nil, ErrNoItems
	}
	var max *Item[E]
	for _, i := range w.items {
		i.actual += i.weight
		if max == nil || i.actual > max.actual {
			max = i
		}
	}
	max.actual -= w.total
	return &max.id, nil
}

func (w *wrrSimple[E]) AddItem(items ...Item[E]) error {
	w.Lock()
	defer w.Unlock()

	for _, i := range items {
		w.items = append(w.items, &i)
	}

	w.reset()
	return nil
}

// Add add a new item to the balancer, this will automatically lead into a reset of the history
func (w *wrrSimple[E]) Add(id E, weight int) error {
	return w.AddItem(Item[E]{
		id:     id,
		weight: weight,
		actual: 0,
	})
}

func (w *wrrSimple[E]) reset() {
	w.total = 0
	for _, i := range w.items {
		i.actual = 0
		w.total += i.weight
	}
}
