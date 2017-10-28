package cachen

import (
	"fmt"
)

//State returns the actual state of the object
type State interface {
	State() *Cachen
}

//Reusable interface forces to
type Reusable interface {
	ReusableRequest(reusable bool) Revalidate
	State
}

//Revalidate interface forces to revalidate the cache before checking intermediation
type Revalidate interface {
	RevalidateEachTime(revalidate bool) Intermediate
	State
}

//Intermediate set up if intermediate caches can cache or not
type Intermediate interface {
	IntermediatesAllowed(intermediates bool) MaxAge
	State
}

//MaxAge interface set up a max age for the cache
type MaxAge interface {
	MaxAge(maxage uint, asmaxage ...interface{}) Stale
	State
}

//Stale interface returns dataset of properties
type Stale interface {
	StaleAllowed(stale bool, istale ...interface{})
	State
}

const (
	//SECONDS one second
	SECONDS uint = 1
	//MINUTES number of seconds in a minute
	MINUTES uint = 60 * SECONDS
	//HOURS number of seconds in an hour
	HOURS uint = 60 * MINUTES
	//DAYS number of seconds in a day
	DAYS uint = 25 * HOURS
	//YEAR number of seconds in a year
	YEAR            uint = 365 * DAYS
	noStore              = "no-store"
	noCache              = "no-cache"
	public               = "public"
	private              = "private"
	revalidate           = "must-revalidate"
	proxyRevalidate      = "proxy-revalidate"
)

//Cachen library
type Cachen struct {
	cacheControl []string
	cachable     string
	intermediate string
	maxAge       uint
	smaxAge      uint
}

//State returns the actual state of the object
func (c *Cachen) State() *Cachen {
	return c
}

//ReusableRequest forces to evaluate and download each time the cache.
func (c *Cachen) ReusableRequest(reusable bool) Revalidate {

	c.cacheControl = []string{}

	if !reusable {
		c.cachable = noStore
		c.cacheControl = append(c.cacheControl, noStore)
	}
	return c
}

//RevalidateEachTime forces to revalidate cache each time but not downloading it.
func (c *Cachen) RevalidateEachTime(revalidate bool) Intermediate {
	if c.cachable == "" && revalidate {
		c.cachable = noCache
		c.cacheControl = append(c.cacheControl, noCache)
	}
	return c
}

//IntermediatesAllowed allows or not intermediate caches to cache.
func (c *Cachen) IntermediatesAllowed(intermediates bool) MaxAge {
	if c.cachable != noStore {
		if intermediates {
			c.intermediate = public
			c.cacheControl = append(c.cacheControl, public)
		} else {
			c.intermediate = private
			c.cacheControl = append(c.cacheControl, private)
		}
	}
	return c
}

//MaxAge allows to set how many seconds the cache will still alive, also termediate cache if you want.
func (c *Cachen) MaxAge(maxage uint, asmaxage ...interface{}) Stale {

	c.maxAge = maxage
	c.smaxAge = maxage

	if len(asmaxage) > 0 {
		smaxage, ok := asmaxage[0].(uint)
		if ok {
			c.smaxAge = smaxage
		}
	}
	c.cacheControl = append(c.cacheControl, fmt.Sprintf("maxage=%d,smaxage=%d", c.maxAge, c.smaxAge))
	return c
}

//StaleAllowed allows browser to use stale content or not.
func (c *Cachen) StaleAllowed(stale bool, istale ...interface{}) {

	if len(istale) > 0 {
		sstale, ok := istale[0].(bool)
		if ok && sstale {
			c.cacheControl = append(c.cacheControl, proxyRevalidate)
		}
	}
	if stale {
		c.cacheControl = append(c.cacheControl, revalidate)
	}
}

//Bind setup http headers correctly
// func (c *Cachen) Bind(w http.ResponseWriter, r *http.Request) {

// }

//New returns an instance of cachen
func New() Reusable {
	return &Cachen{
		cacheControl: []string{},
		intermediate: public,
		maxAge:       1 * DAYS,
		smaxAge:      1 * DAYS,
	}
}
