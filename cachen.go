package cachen

import "net/http"

//Reusable interface returns the function to setup if the cache is or not reusable
type Reusable interface {
	ReusableRequest(reusable bool)
}

const (
	noStore = "no-store"
	noCache = "no-cache"
)

type cachen struct {
	cachable     string
	intermediate int16
}

//ReusableRequest evaluates if the actual request is or not reusable
func (c *cachen) ReusableRequest(reusable bool) {
	if !reusable {
		c.cachable = noStore
	}
}

//RevalidateEachTime evaluates if the actual request is or not reusable
func (c *cachen) RevalidateEachTime(revalidate bool) *cachen {
	if c.cachable == "" && revalidate {
		c.cachable = noCache
	}
	return c
}

//Handler executes the configuration
func (c *cachen) Handler(w http.ResponseWriter, r *http.Request) {

}

//New returns an instance of cachen
func New() Reusable {
	return &cachen{}
}
