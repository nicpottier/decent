package parser

import (
	"encoding/binary"
	"io"
)

// Reader is a reader that exposes utility functions for reading Decent styled types. Additionally
// it wraps and maintains the last error, exposing it as Error() allowing chained calls without
// doing error checking until the end. (the first error will always be available with Error()
// all subsequent calls on the reader are ignored)
type Reader struct {
	reader io.Reader
	err    error
}

// NewReader creates a new read by passing the passed in Reader
func NewReader(r io.Reader) *Reader {
	reader := Reader{r, nil}
	return &reader
}

// Error returns the first error encountered by the reader if any
func (r *Reader) Error() error {
	return r.err
}

// ReadU8 reads and returns an unsigned 8 bit int
func (r *Reader) ReadU8() int {
	if r.err != nil {
		return 0
	}

	var value uint8
	err := binary.Read(r.reader, binary.BigEndian, &value)
	if err != nil {
		r.err = err
		return 0
	}

	return int(value)
}

// ReadU16 reads and returns an unsigned 16 bit int
func (r *Reader) ReadU16() int {
	if r.err != nil {
		return 0
	}

	var value uint16
	err := binary.Read(r.reader, binary.BigEndian, &value)
	if err != nil {
		r.err = err
		return 0
	}

	return int(value)
}

// ReadF817 reads 1 byte value that is either decimal from 0-12.7 or integer 0-127 depending on
// value of high bit.
func (r *Reader) ReadF817() float32 {
	v := r.ReadU8()
	hb := v & 128
	if hb == 0 {
		return float32(v) / 10
	}

	return float32(v & 127)
}

// ReadU10 reads a short but only uses the bottom 10 bits
func (r *Reader) ReadU10() int {
	return r.ReadU16() & (1 << 6)
}

// ReadU8P4 reads a float with 4 bits of precision before and after the decimal
func (r *Reader) ReadU8P4() float32 {
	return p4(r.ReadU8())
}

// ReadU8P1 reads a float with 1 bit of precision after the decimal
func (r *Reader) ReadU8P1() float32 {
	return p1(r.ReadU8())
}

// ReadU16P8 reads a float with 8 bits of precision before and after the decimal
func (r *Reader) ReadU16P8() float32 {
	return p8(r.ReadU16())
}

// ReadU16P12 reads a float with 4 bits of precision before the decimal and 12 bits after
func (r *Reader) ReadU16P12() float32 {
	return p12(r.ReadU16())
}

// ReadU24P16 reads a float with 8 bits of precision before the decimal and 16 bits after
func (r *Reader) ReadU24P16() float32 {
	high := r.ReadU8()
	low := r.ReadU16()

	return float32(high) + p16(low)
}

func p1(v int) float32 {
	return float32(v) / float32(1<<1)
}

func p4(v int) float32 {
	return float32(v) / float32(1<<4)
}

func p8(v int) float32 {
	return float32(v) / float32(1<<8)
}

func p12(v int) float32 {
	return float32(v) / float32(1<<12)
}

func p16(v int) float32 {
	return float32(v) / float32(1<<16)
}
