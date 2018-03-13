package app

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/cristaloleg/bc/block"
)

func (a *App) home(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadFile("app/page.html")
	w.Write(data)
}

func (a *App) blocks(w http.ResponseWriter, r *http.Request) {
	it := a.bc.Iterator()
	data := []*block.Block{}

	for it.HasNext() {
		if b := it.Next(); b != nil {
			data = append(data, b)
		}
	}

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *App) mineBlock(w http.ResponseWriter, r *http.Request) {
	m, _ := url.ParseQuery(r.URL.RawQuery)
	data, ok := m["data"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hash, err := a.bc.AddBlock([]byte(data[0]))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	a.hub.Broadcast("mined block: " + hex.EncodeToString(hash))

	newBlock := struct {
		Hash string
	}{
		Hash: hex.EncodeToString(hash),
	}
	json.NewEncoder(w).Encode(newBlock)
}

func (a *App) peers(w http.ResponseWriter, r *http.Request) {
	peers := a.hub.GetPeers()
	data := []string{}

	for _, p := range peers {
		if p != nil {
			data = append(data, p.conn.RemoteAddr().String())
		}
	}

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *App) ws(w http.ResponseWriter, r *http.Request) {
	c, err := a.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	client := &wsClient{
		hub:  a.hub,
		conn: c,
		send: make(chan []byte, 10),
	}

	a.hub.Connect(client)

	go client.writePump()
	go client.readPump()
}
