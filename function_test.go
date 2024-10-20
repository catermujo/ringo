package ringo

import (
	"testing"

	. "gopkg.in/check.v1"
)

// hook up go-check to go testing
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestFindPowerOfTwo(c *C) {
	// given
	cap1 := uint64(0)
	cap2 := uint64(10)
	cap3 := uint64(16)
	cap4 := uint64(33)
	cap5 := ^uint64(0)

	// when
	res1 := findPowerOfTwo(cap1)
	res2 := findPowerOfTwo(cap2)
	res3 := findPowerOfTwo(cap3)
	res4 := findPowerOfTwo(cap4)
	res5 := findPowerOfTwo(cap5)

	// then
	c.Assert(res1, Equals, uint64(0))
	c.Assert(res2, Equals, uint64(16))
	c.Assert(res3, Equals, uint64(16))
	c.Assert(res4, Equals, uint64(64))
	c.Assert(res5, Equals, uint64(0))
}

var bufferSet = []BufferType{NodeBased, Classical}

func (s *MySuite) TestOfferAndPollSuccess(c *C) {
	for _, t := range bufferSet {
		// given
		fakeString := "fake"
		buffer := New[*string](t, 10)

		// when
		result := buffer.Put(&fakeString)
		poll, _ := buffer.Get()

		// then
		c.Assert(result, Equals, true)
		c.Assert(poll, Equals, &fakeString)
	}
}

func (s *MySuite) TestOfferFailedWhenFull(c *C) {
	for _, t := range bufferSet {
		// given
		capacity := 10
		buffer := New[int](t, uint64(capacity))
		realCapacity := findPowerOfTwo(uint64(capacity + 1))
		for i := 0; i < int(realCapacity); i++ {
			buffer.Put(i)
		}

		// when
		offered := buffer.Put(10)

		// then
		c.Assert(offered, Equals, false)
	}
}

func (s *MySuite) TestPollFailedWhenEmpty(c *C) {
	for _, t := range bufferSet {
		// given
		capacity := 10
		buffer := New[int](t, uint64(capacity))

		// when
		_, success := buffer.Get()

		// then
		c.Assert(success, Equals, false)
	}
}

func (s *MySuite) TestRingBufferShift(c *C) {
	for _, t := range bufferSet {
		// given
		capacity := 10
		buffer := New[int](t, uint64(capacity))

		// when
		for i := 0; i < 13; i++ {
			buffer.Put(i)
		}

		// when
		buffer.Put(13)
		buffer.Put(14)

		// then
		polled, success := buffer.Get()
		c.Assert(success, Equals, true)
		c.Assert(polled, Equals, 0)

		// when
		buffer.Put(15)

		// then
		for i := 0; i < 14; i++ {
			polled, success := buffer.Get()
			c.Assert(success, Equals, true)
			c.Assert(polled, Equals, i+1)
		}

		// when
		buffer.Put(16)
		buffer.Put(17)
		buffer.Put(18)

		// then
		polled1, _ := buffer.Get()
		c.Assert(polled1, Equals, 15)
		polled2, _ := buffer.Get()
		c.Assert(polled2, Equals, 16)
	}
}
