package client

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
)

type Decoding struct {
	ID    int             `json:"id"`
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}

func decoding(data []byte) *Decoding {
	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		fmt.Println("Error creating zlib reader:", err)
		return &Decoding{}
	}
	defer reader.Close()

	var decompressedData bytes.Buffer
	_, err = decompressedData.ReadFrom(reader)
	if err != nil {
		fmt.Println("Error decompressing message:", err)
		return &Decoding{}
	}

	var msg Decoding
	err = json.Unmarshal(decompressedData.Bytes(), &msg)
	if err != nil {
		fmt.Println("Error decoding message:", err)
		return &Decoding{}
	}

	return &msg
}
