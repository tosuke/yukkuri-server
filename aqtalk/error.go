package aqtalk

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidKoeType = errors.New("invalid koe type")
)

type AqTalk1Error struct {
	Code int
}

func (aqtkerr *AqTalk1Error) Error() string {
	switch aqtkerr.Code {
	case 100: return "misc error"
	case 101: return "insufficient memory"
	case 102: case 105: return "invalid symbols"
	case 103: return "negative prosody data"
	case 104: return "internal error"
	case 106: return "invalid tag"
	case 107: return "too long tag"
	case 108: return "invalid tag value"
	case 109: return "cannot play wave(sound driver)"
	case 110: return "cannot play wave asynchronously(sound driver)"
	case 111: return "no data to synthesize"
	case 200: case 202: return "too long koe"
	case 201: return "too many symbols in a phrase"
	case 203: return "heap memory exhaust"
	case 205: return "license key is needed"
	default: return fmt.Sprintf("invalid koe(pos:%d)", aqtkerr.Code)
	}
	panic("unreachable")
}