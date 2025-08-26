# WeightedRoundRobinSimple
A simple implementation of a weighted round robin (WRR) algorithm.

## Install 

simply do

`go get github.com/willie68/gowillie68/WeightedRoundRobinSimple`

## Usage

```go
func main() {
    wrr := New[string]()

    wrr.Add("1", 10)
	wrr.Add("2", 5)
	wrr.Add("3", 2)

    id, _ := wrr.GetNext()
    fmt.Printf("id is %s", id)
}
```

