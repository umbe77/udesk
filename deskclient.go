package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
)

func convertPathToUrl(imagePath string) string {
	if imagePath == "" {
		return imagePath
	}
	img, err := os.ReadFile(imagePath)
	if err != nil {
		return err.Error()
	}
	mimeType := http.DetectContentType(img)
	if filepath.Ext(imagePath) == ".svg" {
		mimeType = "image/svg+xml"
	}
	encoded := base64.StdEncoding.EncodeToString(img)

	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded)
}

type Client struct {
	selectedActions map[string]SelectAction
}

func NewClient() *Client {
	return &Client{
		selectedActions: map[string]SelectAction{
			"processes":    killProc,
			"applications": runProc,
		},
	}
}

func runProc(ctx context.Context, item ListItem) {
	cmd := item.Original["execPath"].(string)
	msg := Message{
		Action: "runapp",
		Params: []string{
			cmd,
		},
	}

	send(msg)
}

func killProc(ctx context.Context, item ListItem) {
	var pid = item.Original["pid"].(float64)
	msg := Message{
		Action: "killproc",
		Params: []string{
			fmt.Sprintf("%.0f", pid),
		},
	}

	send(msg)
}

func send(msg Message) (ResponseMessage, error) {
	addr, err := xdg.DataFile("deskservice/action.sock")
	if err != nil {
		return ResponseMessage{
			Action: msg.Action,
			Code:   int8(ERR),
			Params: []string{
				err.Error(),
			},
		}, err
	}

	conn, err := net.Dial("unix", addr)
	if err != nil {
		return ResponseMessage{
			Action: msg.Action,
			Code:   int8(ERR),
			Params: []string{
				err.Error(),
			},
		}, err
	}
	defer conn.Close()

	m, err := Encode(msg)
	if err != nil {
		return ResponseMessage{
			Action: msg.Action,
			Code:   int8(ERR),
			Params: []string{
				err.Error(),
			},
		}, err
	}

	if _, err := conn.Write(m); err != nil {
		return ResponseMessage{
			Action: msg.Action,
			Code:   int8(ERR),
			Params: []string{
				err.Error(),
			},
		}, err
	}

	reader := bufio.NewReader(conn)
	respMsg, err := reader.ReadBytes('\n')
	if err != nil {
		if err == io.EOF {
			return ResponseMessage{
				Action: msg.Action,
				Code:   int8(ERR),
				Params: []string{
					err.Error(),
				},
			}, err
		}
		return ResponseMessage{
			Action: msg.Action,
			Code:   int8(ERR),
			Params: []string{
				err.Error(),
			},
		}, err
	}
	respMsg = respMsg[:len(respMsg)-1]

	message, err := DecodeResponse(respMsg)
	if err != nil {
		return ResponseMessage{
			Action: msg.Action,
			Code:   int8(ERR),
			Params: []string{
				err.Error(),
			},
		}, err
	}

	return message, nil
}

func (c *Client) GetProcesses() ([]ListItem, error) {
	msg := Message{
		Action: "processes",
		Params: make([]string, 0),
	}
	resp, err := send(msg)
	if err != nil {
		return []ListItem{}, nil
	}

	result := make([]ListItem, 0)

	uq := make(map[float64]ListItem)

	for _, item := range resp.Params {
		var lItem map[string]any
		if err := json.Unmarshal([]byte(item), &lItem); err != nil {
			return []ListItem{}, err
		}
		pid := lItem["pid"].(float64)
		if _, ok := uq[pid]; !ok {
			uq[pid] = ListItem{
				Icon:        "",
				Text:        lItem["cmd"].(string),
				Description: item,
				Original:    lItem,
			}
			result = append(result, uq[pid])
		}
	}

	return result, nil
}

func (c *Client) GetApplications() ([]ListItem, error) {
	msg := Message{
		Action: "applications",
		Params: make([]string, 0),
	}
	resp, err := send(msg)
	if err != nil {
		return []ListItem{}, err
	}

	result := make([]ListItem, len(resp.Params))
	for i, item := range resp.Params {
		var lItem map[string]any
		if err := json.Unmarshal([]byte(item), &lItem); err != nil {
			return []ListItem{}, err
		}
		icon := convertPathToUrl(strings.TrimSpace(lItem["iconPath"].(string)))
		result[i] = ListItem{
			Icon:        icon,
			Text:        lItem["name"].(string),
			Description: lItem["execPath"].(string),
			Original:    lItem,
		}

	}

	return result, nil
}
