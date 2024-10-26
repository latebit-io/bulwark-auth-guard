package bulwark

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func doPut(ctx context.Context, url string, payload interface{}, client *http.Client) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		jsonError := &JsonError{}

		if err := json.NewDecoder(resp.Body).Decode(jsonError); err != nil {
			return err
		}

		if jsonError != nil {
			return fmt.Errorf("%s - %s", jsonError.Title, jsonError.Detail)
		}
	}

	return nil
}
