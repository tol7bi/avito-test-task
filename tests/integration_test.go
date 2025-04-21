package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

const baseURL = "http://localhost:8080"

func postJSON(path string, token string, body any) ([]byte, error) {
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", baseURL+path, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func getToken(role string) string {
	raw, _ := postJSON("/dummyLogin", "", map[string]string{"role": role})
	var token string
	_ = json.Unmarshal(raw, &token)
	return token
}

func TestIntegrationPVZFlow(t *testing.T) {
	modToken := getToken("moderator")
	empToken := getToken("employee")

	// 1. Создать ПВЗ
	pvzRaw, _ := postJSON("/pvz", modToken, map[string]string{"city": "Казань"})
	var pvz struct {
		ID string `json:"id"`
	}
	json.Unmarshal(pvzRaw, &pvz)

	// 2. Создать приёмку
	_, err := postJSON("/receptions", empToken, map[string]string{"pvzId": pvz.ID})
	if err != nil {
		t.Fatal("failed to create reception:", err)
	}

	// 3. Добавить 50 товаров
	for i := 0; i < 50; i++ {
		_, err := postJSON("/products", empToken, map[string]string{
			"type":  "электроника",
			"pvzId": pvz.ID,
		})
		if err != nil {
			t.Fatalf("failed to add product #%d: %v", i, err)
		}
	}

	// 4. Закрыть приёмку
	closeResp, err := postJSON(fmt.Sprintf("/pvz/%s/close_last_reception", pvz.ID), empToken, nil)
	if err != nil {
		t.Fatal("failed to close reception:", err)
	}
	var closed map[string]any
	_ = json.Unmarshal(closeResp, &closed)
	if closed["status"] != "close" {
		t.Fatalf("reception not closed successfully: %+v", closed)
	}

	// 5. Проверить GET /pvz
	req, _ := http.NewRequest("GET", baseURL+"/pvz", nil)
	req.Header.Set("Authorization", "Bearer "+empToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal("GET /pvz failed:", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var result []struct {
		PVZ struct {
			ID string `json:"id"`
		} `json:"pvz"`
		Receptions []struct {
			Products []any `json:"products"`
		} `json:"receptions"`
	}
	_ = json.Unmarshal(body, &result)

	var found bool
	for _, item := range result {
		if item.PVZ.ID == pvz.ID {
			found = true
			if len(item.Receptions) != 1 {
				t.Fatalf("expected 1 reception, got %d", len(item.Receptions))
			}
			if len(item.Receptions[0].Products) != 50 {
				t.Fatalf("expected 50 products, got %d", len(item.Receptions[0].Products))
			}
		}
	}
	if !found {
		t.Fatal("created PVZ not found in GET /pvz response")
	}
}
