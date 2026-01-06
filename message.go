package main

import (
	"context"

	"github.com/shamaton/msgpack/v2"
)

type ResponseCode int8

const (
	OK ResponseCode = iota
	NOACTION
	INVALIDNUMOFPARAMS
	ERR
)

type ResponseMessage struct {
	Code   int8
	Action string
	Params []string
}

type Message struct {
	Action string
	Params []string
}

type ListItem struct {
	Icon        string         `json:"icon"`
	Text        string         `json:"text"`
	Description string         `json:"description"`
	Original    map[string]any `json:"original"`
}

var Actions = []string{
	"processes",
	"applications",
}

type SelectAction func(ctx context.Context, item ListItem)

func DecodeResponse(msg []byte) (ResponseMessage, error) {
	var m ResponseMessage
	err := msgpack.Unmarshal(msg, &m)
	if err != nil {
		return ResponseMessage{}, err
	}

	return m, nil
}

func Encode(msg Message) ([]byte, error) {
	resp, err := msgpack.Marshal(msg)
	if err != nil {
		return nil, err
	}
	resp = append(resp, byte('\n'))
	return resp, nil
}
