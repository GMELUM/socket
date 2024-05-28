package connect

import (
	"bytes"

	"compress/zlib"
	"encoding/json"
)

type Decoding struct {
	ID    int             `json:"id"`
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}

func (r *Decoding) UnmarshalJSON(p []byte) error {
	var tmp []json.RawMessage
	if err := json.Unmarshal(p, &tmp); err != nil {
		return err
	}
	if err := json.Unmarshal(tmp[0], &r.ID); err != nil {
		return err
	}
	if err := json.Unmarshal(tmp[1], &r.Type); err != nil {
		return err
	}
	if err := json.Unmarshal(tmp[2], &r.Value); err != nil {
		return err
	}
	return nil
}

func decoding(data []byte) (*Decoding, error) {
	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var decompressedData bytes.Buffer
	_, err = decompressedData.ReadFrom(reader)
	if err != nil {
		return nil, err
	}

	var msg Decoding
	err = json.Unmarshal(decompressedData.Bytes(), &msg)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}
