package connect

import (
	"bytes"
	"fmt"

	"compress/zlib"
	"encoding/json"
)

type Encoding struct {
	ID    int         `json:"id"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

func encoding(msg Encoding) *[]byte {
	jsonData, err := json.Marshal([]interface{}{msg.ID, msg.Type, msg.Value})
	if err != nil {
		fmt.Println("Error encoding message:", err)
		return nil
	}

	var compressedData bytes.Buffer
	writer := zlib.NewWriter(&compressedData)
	_, err = writer.Write(jsonData)
	if err != nil {
		fmt.Println("Error compressing message:", err)
		return nil
	}
	writer.Close()

	data := compressedData.Bytes()
	return &data
}
