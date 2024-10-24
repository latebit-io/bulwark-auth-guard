package bulwark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func doPost(url string, payload interface{}, model interface{}, client *http.Client) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	if resp.StatusCode >= 300 {
		jsonError := &JsonError{}
		if err := json.NewDecoder(resp.Body).Decode(jsonError); err != nil {
			return err
		}
		if jsonError != nil {
			return fmt.Errorf("%s - %s", jsonError.Title, jsonError.Detail)
		}
	}

	if resp.Body != http.NoBody {
		if err := json.NewDecoder(resp.Body).Decode(model); err != nil {
			return err
		}
	}

	defer resp.Body.Close()

	return nil
}
