package bybit

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Core data types
type Holding struct {
	Coin   string `json:"coin"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

// Price data structures
type TickerPriceResponse struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		Category string `json:"category"`
		List     []struct {
			Symbol    string `json:"symbol"`
			LastPrice string `json:"lastPrice"`
		} `json:"list"`
	} `json:"result"`
}

// HTTP client for requests
var httpClient = &http.Client{Timeout: 4 * time.Second}

// firstNonEmpty returns the first argument that renders to a non-empty string
func firstNonEmpty(values ...interface{}) string {
	for _, v := range values {
		if v == nil {
			continue
		}
		s := fmt.Sprintf("%v", v)
		if s != "" && s != "<nil>" {
			return s
		}
	}
	return ""
}

// signV5 creates the Bybit v5 HMAC SHA256 signature
func signV5(apiKey, secret string, query url.Values, timestamp string, recvWindow string) string {
	// v5 signature payload: timestamp + apiKey + recvWindow + queryString (only the actual URL query)
	// Build canonical query matching what we will send
	var keys []string
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var encoded string
	for i, k := range keys {
		if i > 0 {
			encoded += "&"
		}
		// Values are simple; use first value
		encoded += fmt.Sprintf("%s=%s", k, query.Get(k))
	}
	payload := timestamp + apiKey + recvWindow + encoded
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

// toNumString converts interface{} to a numeric-like string, defaulting nil/"<nil>"/"" to "0"
func toNumString(v interface{}) string {
	if v == nil {
		return "0"
	}
	switch t := v.(type) {
	case string:
		if t == "" || t == "<nil>" {
			return "0"
		}
		return t
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.FormatInt(t, 10)
	default:
		s := fmt.Sprintf("%v", v)
		if s == "<nil>" || s == "" {
			return "0"
		}
		return s
	}
}

// getSpotHoldings fetches wallet balances via REST v5
func (s *BybitService) getSpotHoldings(ctx context.Context, userID string) ([]Holding, error) {
	creds, err := s.GetBybitByUserId(userID)
	if err != nil {
		return nil, fmt.Errorf("bybit credentials not found: %w", err)
	}

	endpoint := "https://api.bybit.com/v5/account/wallet-balance"
	accountType := "UNIFIED"
	recvWindow := "8000"
	timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())

	// Actual query string
	q := url.Values{}
	q.Set("accountType", accountType)

	signature := signV5(creds.ApiKey, creds.ApiSecret, q, timestamp, recvWindow)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-BAPI-API-KEY", creds.ApiKey)
	req.Header.Set("X-BAPI-TIMESTAMP", timestamp)
	req.Header.Set("X-BAPI-RECV-WINDOW", recvWindow)
	req.Header.Set("X-BAPI-SIGN", signature)
	req.Header.Set("X-BAPI-SIGN-TYPE", "2")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("bybit wallet-balance error: %s", string(body))
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	// Handle Bybit business error even with HTTP 200
	if rc, ok := raw["retCode"].(float64); ok && rc != 0 {
		retMsg := fmt.Sprintf("%v", raw["retMsg"])
		return nil, fmt.Errorf("bybit error retCode=%v: %s", rc, retMsg)
	}

	var holdings []Holding
	result, _ := raw["result"].(map[string]interface{})
	if result != nil {
		list, _ := result["list"].([]interface{})
		for _, item := range list {
			iMap, _ := item.(map[string]interface{})
			coins, _ := iMap["coin"].([]interface{})
			for _, c := range coins {
				cm, _ := c.(map[string]interface{})
				coin := toNumString(cm["coin"]) // symbol string
				if coin == "" {
					continue
				}
				// Prefer total wallet balance if present
				total := toNumString(cm["walletBalance"]) // unified total
				if total == "0" || total == "" {
					// Fallback to free+locked if provided
					free := toNumString(cm["free"])     // may be nil
					locked := toNumString(cm["locked"]) // may be nil
					// store separately; UI will sum
					holdings = append(holdings, Holding{Coin: coin, Free: free, Locked: locked})
				} else {
					// Put total into Free to represent overall amount; Locked set to 0
					holdings = append(holdings, Holding{Coin: coin, Free: total, Locked: "0"})
				}
			}
		}
	}
	return holdings, nil
}

// getCurrentPrice gets current price for a symbol via REST API
func getCurrentPrice(symbol string) (string, error) {
	// Format symbol for Bybit API (e.g., "btc" -> "BTCUSDT")
	symbol = fmt.Sprintf("%sUSDT", strings.ToUpper(symbol))

	// Build request URL
	baseURL := "https://api.bybit.com/v5/market/tickers"
	params := url.Values{}
	params.Set("category", "spot")
	params.Set("symbol", symbol)

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Make request
	resp, err := httpClient.Get(reqURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch price: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// Parse response
	var priceResp TickerPriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&priceResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if priceResp.RetCode != 0 {
		return "", fmt.Errorf("API error: %s", priceResp.RetMsg)
	}

	if len(priceResp.Result.List) == 0 {
		return "", fmt.Errorf("no price data found for symbol %s", symbol)
	}

	return priceResp.Result.List[0].LastPrice, nil
}
