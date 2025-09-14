package bybit

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type PriceData struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
	Time   int64  `json:"time"`
}

type WebSocketManager struct {
	conn          *websocket.Conn
	subscribers   map[string][]chan PriceData
	mu            sync.RWMutex
	isConnected   bool
	reconnectChan chan bool
	ctx           context.Context
	cancel        context.CancelFunc
}

type TickerResponse struct {
	Topic string `json:"topic"`
	Type  string `json:"type"`
	Data  struct {
		Symbol    string `json:"symbol"`
		LastPrice string `json:"lastPrice"`
		Ts        int64  `json:"ts"`
	} `json:"data"`
}

var wsManager *WebSocketManager
var wsOnce sync.Once

func GetWebSocketManager() *WebSocketManager {
	wsOnce.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		wsManager = &WebSocketManager{
			subscribers:   make(map[string][]chan PriceData),
			reconnectChan: make(chan bool, 1),
			ctx:           ctx,
			cancel:        cancel,
		}
		go wsManager.connectionManager()
	})
	return wsManager
}

func (ws *WebSocketManager) connectionManager() {
	for {
		select {
		case <-ws.ctx.Done():
			return
		default:
			if err := ws.connect(); err != nil {
				log.Printf("WebSocket connection failed: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}
			ws.handleMessages()
		}
	}
}

func (ws *WebSocketManager) connect() error {
	log.Printf("Attempting to connect to Bybit WebSocket...")
	conn, _, err := websocket.DefaultDialer.Dial("wss://stream.bybit.com/v5/public/spot", nil)
	if err != nil {
		return err
	}

	ws.mu.Lock()
	ws.conn = conn
	ws.isConnected = true
	ws.mu.Unlock()

	log.Printf("WebSocket connected successfully")
	return nil
}

func (ws *WebSocketManager) handleMessages() {
	defer func() {
		ws.mu.Lock()
		if ws.conn != nil {
			ws.conn.Close()
		}
		ws.isConnected = false
		ws.mu.Unlock()
	}()

	for {
		select {
		case <-ws.ctx.Done():
			return
		default:
			var response TickerResponse
			if err := ws.conn.ReadJSON(&response); err != nil {
				log.Printf("WebSocket read error: %v", err)
				return
			}

			if response.Type == "snapshot" || response.Type == "delta" {
				priceData := PriceData{
					Symbol: response.Data.Symbol,
					Price:  response.Data.LastPrice,
					Time:   response.Data.Ts,
				}
				ws.broadcast(priceData)
			}
		}
	}
}

func (ws *WebSocketManager) broadcast(data PriceData) {
	ws.mu.RLock()
	subscribers, exists := ws.subscribers[data.Symbol]
	ws.mu.RUnlock()

	if exists {
		for _, ch := range subscribers {
			select {
			case ch <- data:
			default:
				// Channel is full, skip
			}
		}
	}
}

func (ws *WebSocketManager) Subscribe(symbol string) (chan PriceData, error) {
	symbol = strings.ToUpper(symbol) + "USDT"

	ch := make(chan PriceData, 10)

	ws.mu.Lock()
	ws.subscribers[symbol] = append(ws.subscribers[symbol], ch)
	isConnected := ws.isConnected
	ws.mu.Unlock()

	// Wait for connection if not connected
	if !isConnected {
		log.Printf("WebSocket not connected, waiting...")
		time.Sleep(2 * time.Second) // Give it some time to connect

		ws.mu.RLock()
		isConnected = ws.isConnected
		ws.mu.RUnlock()
	}

	if isConnected {
		if err := ws.subscribeToSymbol(symbol); err != nil {
			log.Printf("Failed to subscribe to symbol %s: %v", symbol, err)
			return nil, err
		}
		log.Printf("Successfully subscribed to %s", symbol)
	} else {
		log.Printf("WebSocket still not connected for symbol %s", symbol)
	}

	return ch, nil
}

func (ws *WebSocketManager) Unsubscribe(symbol string, ch chan PriceData) {
	symbol = strings.ToUpper(symbol) + "USDT"

	ws.mu.Lock()
	defer ws.mu.Unlock()

	subscribers := ws.subscribers[symbol]
	for i, subscriber := range subscribers {
		if subscriber == ch {
			close(ch)
			ws.subscribers[symbol] = append(subscribers[:i], subscribers[i+1:]...)
			break
		}
	}

	if len(ws.subscribers[symbol]) == 0 {
		delete(ws.subscribers, symbol)
		if ws.isConnected {
			ws.unsubscribeFromSymbol(symbol)
		}
	}
}

func (ws *WebSocketManager) subscribeToSymbol(symbol string) error {
	if ws.conn == nil {
		return fmt.Errorf("connection not established")
	}

	subscribeMsg := map[string]interface{}{
		"op":   "subscribe",
		"args": []string{"tickers." + symbol},
	}

	return ws.conn.WriteJSON(subscribeMsg)
}

func (ws *WebSocketManager) unsubscribeFromSymbol(symbol string) error {
	if ws.conn == nil {
		return fmt.Errorf("connection not established")
	}

	unsubscribeMsg := map[string]interface{}{
		"op":   "unsubscribe",
		"args": []string{"tickers." + symbol},
	}

	return ws.conn.WriteJSON(unsubscribeMsg)
}

func (ws *WebSocketManager) Close() {
	if ws.cancel != nil {
		ws.cancel()
	}

	ws.mu.Lock()
	defer ws.mu.Unlock()

	if ws.conn != nil {
		ws.conn.Close()
	}

	for _, subscribers := range ws.subscribers {
		for _, ch := range subscribers {
			close(ch)
		}
	}
	ws.subscribers = make(map[string][]chan PriceData)
}
