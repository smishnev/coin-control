package bybit

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type IconEntry struct {
	Coin string `json:"coin"`
	// Back-compat: keep iconUrl (will mirror darkUrl by default)
	IconURL      string `json:"iconUrl,omitempty"`
	DarkURL      string `json:"darkUrl,omitempty"`
	LightURL     string `json:"lightUrl,omitempty"`
	DarkDataURL  string `json:"darkDataUrl,omitempty"`
	LightDataURL string `json:"lightDataUrl,omitempty"`
}

type iconCacheEntry struct {
	url       string
	fetchedAt time.Time
}

var (
	iconCacheMu sync.Mutex
	iconCache   = map[string]iconCacheEntry{}
	iconTTL     = 24 * time.Hour
	iconIndexTS time.Time
	iconIndex   map[string]iconPair
)

type iconPair struct {
	dark  string
	light string
}

func getCoinIconFromCache(coin string) (string, bool) {
	iconCacheMu.Lock()
	defer iconCacheMu.Unlock()
	e, ok := iconCache[coin]
	if !ok {
		return "", false
	}
	if time.Since(e.fetchedAt) > iconTTL {
		delete(iconCache, coin)
		return "", false
	}
	return e.url, true
}

func setCoinIconCache(coin, urlStr string) {
	iconCacheMu.Lock()
	iconCache[coin] = iconCacheEntry{url: urlStr, fetchedAt: time.Now()}
	iconCacheMu.Unlock()
}

// build or refresh index of coin->iconUrl by querying without coin filter
func ensureIconIndex() error {
	iconCacheMu.Lock()
	if iconIndex != nil && time.Since(iconIndexTS) < iconTTL {
		iconCacheMu.Unlock()
		dbg("Icon index is still fresh, using cached data")
		return nil
	}
	iconCacheMu.Unlock()

	dbg("Building new icon index from Bybit X-API brief-symbol-list...")

	// Try preferred public endpoint that exposes icon URLs (di/li)
	assembled := map[string]iconPair{}
	if idx, err := fetchIconIndexFromBriefList(); err == nil && len(idx) > 0 {
		for k, v := range idx {
			assembled[k] = v
		}
		dbg("Icon index loaded from brief-symbol-list with %d entries", len(idx))
	} else if err != nil {
		dbg("brief-symbol-list fetch failed: %v", err)
	}

	// Fallback: static list for popular coins
	basicIcons := map[string]string{
		"BTC":    "https://s2.coinmarketcap.com/static/img/coins/64x64/1.png",
		"ETH":    "https://s2.coinmarketcap.com/static/img/coins/64x64/1027.png",
		"USDT":   "https://s2.coinmarketcap.com/static/img/coins/64x64/825.png",
		"BNB":    "https://s2.coinmarketcap.com/static/img/coins/64x64/1839.png",
		"XRP":    "https://s2.coinmarketcap.com/static/img/coins/64x64/52.png",
		"ADA":    "https://s2.coinmarketcap.com/static/img/coins/64x64/2010.png",
		"SOL":    "https://s2.coinmarketcap.com/static/img/coins/64x64/5426.png",
		"MATIC":  "https://s2.coinmarketcap.com/static/img/coins/64x64/3890.png",
		"LINK":   "https://s2.coinmarketcap.com/static/img/coins/64x64/1975.png",
		"UNI":    "https://s2.coinmarketcap.com/static/img/coins/64x64/7083.png",
		"AVAX":   "https://s2.coinmarketcap.com/static/img/coins/64x64/5805.png",
		"ATOM":   "https://s2.coinmarketcap.com/static/img/coins/64x64/3794.png",
		"FIL":    "https://s2.coinmarketcap.com/static/img/coins/64x64/2280.png",
		"NEAR":   "https://s2.coinmarketcap.com/static/img/coins/64x64/6535.png",
		"APT":    "https://s2.coinmarketcap.com/static/img/coins/64x64/21794.png",
		"SUI":    "https://s2.coinmarketcap.com/static/img/coins/64x64/20947.png",
		"OP":     "https://s2.coinmarketcap.com/static/img/coins/64x64/11840.png",
		"ARB":    "https://s2.coinmarketcap.com/static/img/coins/64x64/1958.png",
		"STRK":   "https://s2.coinmarketcap.com/static/img/coins/64x64/28744.png",
		"PEPE":   "https://s2.coinmarketcap.com/static/img/coins/64x64/24478.png",
		"DOGE":   "https://s2.coinmarketcap.com/static/img/coins/64x64/74.png",
		"SHIB":   "https://s2.coinmarketcap.com/static/img/coins/64x64/5994.png",
		"LTC":    "https://s2.coinmarketcap.com/static/img/coins/64x64/2.png",
		"BCH":    "https://s2.coinmarketcap.com/static/img/coins/64x64/1831.png",
		"ETC":    "https://s2.coinmarketcap.com/static/img/coins/64x64/1321.png",
		"XLM":    "https://s2.coinmarketcap.com/static/img/coins/64x64/512.png",
		"TRX":    "https://s2.coinmarketcap.com/static/img/coins/64x64/1958.png",
		"CHZ":    "https://s2.coinmarketcap.com/static/img/coins/64x64/4066.png",
		"MANTA":  "https://s2.coinmarketcap.com/static/img/coins/64x64/28736.png",
		"JTO":    "https://s2.coinmarketcap.com/static/img/coins/64x64/28735.png",
		"WLD":    "https://s2.coinmarketcap.com/static/img/coins/64x64/28734.png",
		"PYTH":   "https://s2.coinmarketcap.com/static/img/coins/64x64/28733.png",
		"BLUR":   "https://s2.coinmarketcap.com/static/img/coins/64x64/28732.png",
		"DYDX":   "https://s2.coinmarketcap.com/static/img/coins/64x64/28731.png",
		"LDO":    "https://s2.coinmarketcap.com/static/img/coins/64x64/28730.png",
		"CRV":    "https://s2.coinmarketcap.com/static/img/coins/64x64/28729.png",
		"PENDLE": "https://s2.coinmarketcap.com/static/img/coins/64x64/28728.png",
		"GMX":    "https://s2.coinmarketcap.com/static/img/coins/64x64/28727.png",
		"RDNT":   "https://s2.coinmarketcap.com/static/img/coins/64x64/28726.png",
		"ZRO":    "https://s2.coinmarketcap.com/static/img/coins/64x64/28725.png",
		"AXL":    "https://s2.coinmarketcap.com/static/img/coins/64x64/28724.png",
		"HBAR":   "https://s2.coinmarketcap.com/static/img/coins/64x64/28723.png",
		"MINA":   "https://s2.coinmarketcap.com/static/img/coins/64x64/28722.png",
		"STG":    "https://s2.coinmarketcap.com/static/img/coins/64x64/28721.png",
		"MNT":    "https://s2.coinmarketcap.com/static/img/coins/64x64/28720.png",
		"ARKM":   "https://s2.coinmarketcap.com/static/img/coins/64x64/28719.png",
		"C98":    "https://s2.coinmarketcap.com/static/img/coins/64x64/28718.png",
		"APEX":   "https://s2.coinmarketcap.com/static/img/coins/64x64/28717.png",
		"W":      "https://s2.coinmarketcap.com/static/img/coins/64x64/28716.png",
		"ZKJ":    "https://s2.coinmarketcap.com/static/img/coins/64x64/28715.png",
		"SCR":    "https://s2.coinmarketcap.com/static/img/coins/64x64/28714.png",
		"ENA":    "https://s2.coinmarketcap.com/static/img/coins/64x64/28713.png",
		"POL":    "https://s2.coinmarketcap.com/static/img/coins/64x64/28712.png",
		"TWT":    "https://s2.coinmarketcap.com/static/img/coins/64x64/28711.png",
		"BEAM":   "https://s2.coinmarketcap.com/static/img/coins/64x64/28710.png",
		"HFT":    "https://s2.coinmarketcap.com/static/img/coins/64x64/28709.png",
		"MAVIA":  "https://s2.coinmarketcap.com/static/img/coins/64x64/28708.png",
		"SAND":   "https://s2.coinmarketcap.com/static/img/coins/64x64/28707.png",
		"ACH":    "https://s2.coinmarketcap.com/static/img/coins/64x64/28706.png",
		"TOMI":   "https://s2.coinmarketcap.com/static/img/coins/64x64/28705.png",
		"PORT3":  "https://s2.coinmarketcap.com/static/img/coins/64x64/28704.png",
		"XAI":    "https://s2.coinmarketcap.com/static/img/coins/64x64/28703.png",
		"SQR":    "https://s2.coinmarketcap.com/static/img/coins/64x64/28702.png",
		"SIS":    "https://s2.coinmarketcap.com/static/img/coins/64x64/28701.png",
		"MAGIC":  "https://s2.coinmarketcap.com/static/img/coins/64x64/28700.png",
		"QTUM":   "https://s2.coinmarketcap.com/static/img/coins/64x64/28700.png",
		"GRT":    "https://s2.coinmarketcap.com/static/img/coins/64x64/28700.png",
		"SATS":   "https://s2.coinmarketcap.com/static/img/coins/64x64/28700.png",
		"PAWS":   "https://s2.coinmarketcap.com/static/img/coins/64x64/28700.png",
		"MEME":   "https://s2.coinmarketcap.com/static/img/coins/64x64/28700.png",
		"MEMEFI": "https://s2.coinmarketcap.com/static/img/coins/64x64/28700.png",
		"MCRT":   "https://s2.coinmarketcap.com/static/img/coins/64x64/28700.png",
		"5IRE":   "https://s2.coinmarketcap.com/static/img/coins/64x64/28700.png",
	}

	// Merge fallback for missing coins (use the same URL for both dark/light)
	for k, u := range basicIcons {
		ku := strings.ToUpper(k)
		if _, ok := assembled[ku]; !ok {
			assembled[ku] = iconPair{dark: u, light: u}
		}
	}
	dbg("Using basic icon mapping for %d additional coins", len(basicIcons))

	if len(assembled) == 0 {
		return fmt.Errorf("no icons available from brief list and fallback")
	}
	iconCacheMu.Lock()
	iconIndex = assembled
	iconIndexTS = time.Now()
	iconCacheMu.Unlock()
	return nil
}

// fetchIconIndexFromBriefList queries Bybit X-API brief symbol list and builds coin->iconURL index
func fetchIconIndexFromBriefList() (map[string]iconPair, error) {
	url := "https://bybit.com/x-api/contract/v5/product/brief-symbol-list"
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("brief-symbol-list error: %s", string(body))
	}

	// parse very defensively: result.list [] or list [] or top-level []
	var raw interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	// Extract array of items
	items := []interface{}{}
	if m, ok := raw.(map[string]interface{}); ok {
		if r, ok := m["result"]; ok {
			switch t := r.(type) {
			case map[string]interface{}:
				if l, ok := t["list"].([]interface{}); ok {
					items = l
				}
			case []interface{}:
				items = t
			}
		}
		if len(items) == 0 {
			if l, ok := m["list"].([]interface{}); ok {
				items = l
			}
		}
	} else if arr, ok := raw.([]interface{}); ok {
		items = arr
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("no items in brief-symbol-list")
	}

	idx := make(map[string]iconPair, len(items))
	for _, it := range items {
		m, _ := it.(map[string]interface{})
		if m == nil {
			continue
		}
		coin := extractBaseCoin(m)
		if coin == "" {
			continue
		}
		// Extract dark and light URLs
		dark := firstNonEmpty(m["di"], m["darkIcon"], m["darkIconUrl"], m["dark_icon"], m["dark_icon_url"])
		light := firstNonEmpty(m["li"], m["lightIcon"], m["lightIconUrl"], m["light_icon"], m["light_icon_url"])
		if dark == "" && light == "" {
			// try direct stringify
			if v, ok := m["di"]; ok {
				dark = fmt.Sprintf("%v", v)
			}
			if v, ok := m["li"]; ok {
				light = fmt.Sprintf("%v", v)
			}
		}
		if dark == "" && light == "" {
			continue
		}
		if dark == "" {
			dark = light
		}
		if light == "" {
			light = dark
		}
		cu := strings.ToUpper(coin)
		if _, exists := idx[cu]; !exists {
			idx[cu] = iconPair{dark: dark, light: light}
		}
	}
	return idx, nil
}

// extractBaseCoin tries to get coin code from brief-symbol-list item
func extractBaseCoin(m map[string]interface{}) string {
	if v, ok := m["baseCoin"]; ok {
		return strings.ToUpper(fmt.Sprintf("%v", v))
	}
	// try common alternatives
	for _, k := range []string{"base", "bc", "coin"} {
		if v, ok := m[k]; ok {
			s := strings.ToUpper(fmt.Sprintf("%v", v))
			if s != "" && s != "<nil>" {
				return s
			}
		}
	}
	// derive from symbol and quoteCoin if present
	sym := fmt.Sprintf("%v", m["symbol"])
	if sym == "" || sym == "<nil>" {
		return ""
	}
	if q, ok := m["quoteCoin"]; ok {
		qs := strings.ToUpper(fmt.Sprintf("%v", q))
		if strings.HasSuffix(strings.ToUpper(sym), qs) {
			return strings.TrimSuffix(strings.ToUpper(sym), qs)
		}
	}
	// fallback: strip common suffixes
	us := strings.ToUpper(sym)
	for _, suf := range []string{"USDT", "USDC", "USD", "BTC", "ETH"} {
		if strings.HasSuffix(us, suf) {
			return strings.TrimSuffix(us, suf)
		}
	}
	return ""
}

func fetchCoinIcon(coin string) (string, error) {
	// Try Bybit asset coin query-info (public)
	endpoint := "https://api.bybit.com/v5/asset/coin/query-info"
	q := url.Values{}
	// Some endpoints expect lowercase coin code; try both later
	q.Set("coin", strings.ToUpper(coin))
	resp, err := httpClient.Get(endpoint + "?" + q.Encode())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("coin query-info error: %s", string(body))
	}
	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return "", err
	}
	// Expected raw.result.rows or raw.result.list; also sometimes raw.result is an array
	result, _ := raw["result"].(map[string]interface{})
	var rows []interface{}
	if result != nil {
		rows, _ = result["rows"].([]interface{})
		if len(rows) == 0 {
			rows, _ = result["list"].([]interface{})
		}
	}
	if len(rows) == 0 {
		// maybe result is directly a list
		if alt, ok := raw["result"].([]interface{}); ok {
			rows = alt
		}
	}
	for _, it := range rows {
		m, _ := it.(map[string]interface{})
		c := fmt.Sprintf("%v", m["coin"])
		if strings.EqualFold(c, coin) {
			// possible keys for icon url
			for _, k := range []string{"icon", "iconUrl", "iconURL", "logo", "logoUrl", "logoURL", "logoURI"} {
				if v, ok := m[k]; ok {
					u := fmt.Sprintf("%v", v)
					if u != "" && u != "<nil>" {
						return u, nil
					}
				}
			}
		}
	}
	return "", fmt.Errorf("icon not found for %s", coin)
}

// getCoinIcons returns icon URLs for requested coins using cache
func (s *BybitService) getCoinIcons(coins []string) ([]IconEntry, error) {
	// ensure global index
	_ = ensureIconIndex()
	uniq := map[string]struct{}{}
	for _, c := range coins {
		if c == "" {
			continue
		}
		uniq[strings.ToUpper(c)] = struct{}{}
	}
	out := make([]IconEntry, 0, len(uniq))
	missing := make([]string, 0)
	for c := range uniq {
		if iconIndex != nil {
			if p, ok := iconIndex[c]; ok {
				darkURL := p.dark
				lightURL := p.light
				darkData := cacheAndDataURL(c, "dark", darkURL)
				lightData := cacheAndDataURL(c, "light", lightURL)
				out = append(out, IconEntry{
					Coin:         c,
					IconURL:      darkURL,
					DarkURL:      darkURL,
					LightURL:     lightURL,
					DarkDataURL:  darkData,
					LightDataURL: lightData,
				})
				continue
			}
		}
		if urlStr, ok := getCoinIconFromCache(c); ok {
			data := cacheAndDataURL(c, "dark", urlStr)
			out = append(out, IconEntry{Coin: c, IconURL: urlStr, DarkURL: urlStr, LightURL: urlStr, DarkDataURL: data, LightDataURL: data})
		} else {
			missing = append(missing, c)
		}
	}
	// fetch missing concurrently via query-info
	type res struct{ c, u string }
	ch := make(chan res, len(missing))
	wg := sync.WaitGroup{}
	for _, c := range missing {
		c := c
		wg.Add(1)
		go func() {
			defer wg.Done()
			if u, err := fetchCoinIcon(c); err == nil {
				setCoinIconCache(c, u)
				ch <- res{c: c, u: u}
			}
		}()
	}
	wg.Wait()
	close(ch)
	for r := range ch {
		data := cacheAndDataURL(r.c, "dark", r.u)
		out = append(out, IconEntry{Coin: r.c, IconURL: r.u, DarkURL: r.u, LightURL: r.u, DarkDataURL: data, LightDataURL: data})
	}
	return out, nil
}

// getCoinIconURLs returns only dark/light URLs (no data URLs, no disk IO) for faster initial load
func (s *BybitService) getCoinIconURLs(coins []string) ([]IconEntry, error) {
	_ = ensureIconIndex()
	uniq := map[string]struct{}{}
	for _, c := range coins {
		if c == "" {
			continue
		}
		uniq[strings.ToUpper(c)] = struct{}{}
	}
	out := make([]IconEntry, 0, len(uniq))
	for c := range uniq {
		if iconIndex != nil {
			if p, ok := iconIndex[c]; ok {
				out = append(out, IconEntry{
					Coin:     c,
					IconURL:  p.dark,
					DarkURL:  p.dark,
					LightURL: p.light,
				})
				continue
			}
		}
		if urlStr, ok := getCoinIconFromCache(c); ok {
			out = append(out, IconEntry{Coin: c, IconURL: urlStr, DarkURL: urlStr, LightURL: urlStr})
		} else {
			// best effort: try immediate query-info once
			if u, err := fetchCoinIcon(c); err == nil {
				setCoinIconCache(c, u)
				out = append(out, IconEntry{Coin: c, IconURL: u, DarkURL: u, LightURL: u})
			}
		}
	}
	return out, nil
}

// Prefetch icons to local cache in background with concurrency limit
func (s *BybitService) prefetchCoinIcons(coins []string) {
	_ = ensureIconIndex()
	// build tasks of (coin, flavor, url)
	type task struct{ c, flavor, url string }
	tasks := make([]task, 0, len(coins)*2)
	seen := map[string]struct{}{}
	for _, c := range coins {
		cu := strings.ToUpper(strings.TrimSpace(c))
		if cu == "" {
			continue
		}
		if _, ok := seen[cu]; ok {
			continue
		}
		seen[cu] = struct{}{}
		if iconIndex != nil {
			if p, ok := iconIndex[cu]; ok {
				if p.dark != "" {
					tasks = append(tasks, task{c: cu, flavor: "dark", url: p.dark})
				}
				if p.light != "" {
					tasks = append(tasks, task{c: cu, flavor: "light", url: p.light})
				}
				continue
			}
		}
		if u, ok := getCoinIconFromCache(cu); ok {
			tasks = append(tasks, task{c: cu, flavor: "dark", url: u})
		}
	}
	if len(tasks) == 0 {
		return
	}

	// worker pool
	const workers = 8
	ch := make(chan task)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range ch {
				_ = cacheIconFile(t.c, t.flavor, t.url)
			}
		}()
	}
	for _, t := range tasks {
		ch <- t
	}
	close(ch)
	wg.Wait()
}

// cacheIconFile stores icon to disk cache without base64 encoding step
func cacheIconFile(coinUpper, flavor, urlStr string) error {
	if urlStr == "" || urlStr == "<nil>" {
		return nil
	}
	dir := getIconCacheDir()
	_ = os.MkdirAll(dir, 0o755)
	ext := filepath.Ext(strings.Split(urlStr, "?")[0])
	if ext == "" {
		ext = ".png"
	}
	filename := fmt.Sprintf("%s_%s%s", strings.ToLower(coinUpper), flavor, ext)
	path := filepath.Join(dir, filename)
	if st, err := os.Stat(path); err == nil {
		if time.Since(st.ModTime()) < 7*24*time.Hour { // fresh enough
			return nil
		}
	}
	return downloadFile(urlStr, path)
}

// cacheAndDataURL downloads the icon to local cache (if needed) and returns a base64 data URL
func cacheAndDataURL(coinUpper string, flavor string, urlStr string) string {
	if urlStr == "" || urlStr == "<nil>" {
		return ""
	}
	dir := getIconCacheDir()
	_ = os.MkdirAll(dir, 0o755)
	ext := filepath.Ext(strings.Split(urlStr, "?")[0])
	if ext == "" {
		ext = ".png"
	}
	filename := fmt.Sprintf("%s_%s%s", strings.ToLower(coinUpper), flavor, ext)
	path := filepath.Join(dir, filename)
	// If file recent enough, reuse
	needDownload := true
	if st, err := os.Stat(path); err == nil {
		if time.Since(st.ModTime()) < 7*24*time.Hour {
			needDownload = false
		}
	}
	if needDownload {
		if err := downloadFile(urlStr, path); err != nil {
			dbg("download icon failed for %s: %v", coinUpper, err)
		}
	}
	b, err := os.ReadFile(path)
	if err != nil || len(b) == 0 {
		return ""
	}
	mime := "image/png"
	switch strings.ToLower(ext) {
	case ".svg":
		mime = "image/svg+xml"
	case ".jpg", ".jpeg":
		mime = "image/jpeg"
	case ".webp":
		mime = "image/webp"
	}
	return fmt.Sprintf("data:%s;base64,%s", mime, base64.StdEncoding.EncodeToString(b))
}

func getIconCacheDir() string {
	base, err := os.UserCacheDir()
	if err != nil || base == "" {
		base = os.TempDir()
	}
	return filepath.Join(base, "coin-control", "icons")
}

func downloadFile(urlStr, dest string) error {
	resp, err := httpClient.Get(urlStr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("download failed: %d %s", resp.StatusCode, string(body))
	}
	tmp := dest + ".part"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		_ = os.Remove(tmp)
		return err
	}
	f.Close()
	return os.Rename(tmp, dest)
}
