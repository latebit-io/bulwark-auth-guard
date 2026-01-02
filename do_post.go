package bulwark

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func doPost(ctx context.Context, url string, payload interface{}, model interface{}, client *http.Client) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		jsonError := &JsonError{}
		if err := json.NewDecoder(resp.Body).Decode(jsonError); err != nil {
			return err
		}
		return fmt.Errorf("%s - %s", jsonError.Title, jsonError.Detail)
	}

	if resp.Body != http.NoBody && model != nil {
		if err := json.NewDecoder(resp.Body).Decode(model); err != nil {
			return err
		}
	}

	return nil
}
