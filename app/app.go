package app

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/cristaloleg/bc/block"
	"github.com/cristaloleg/bc/storage"

	"github.com/gorilla/websocket"
)

// App ...
type App struct {
	host       string
	port       string
	bc         *block.Blockchain
	wsUpgrader websocket.Upgrader
	hub        *wsHub
}

// New ...
func New() *App {
	a := &App{
		wsUpgrader: websocket.Upgrader{},
		hub:        newHub(),
	}
	return a
}

// Init ...
func (a *App) Init(configName string) error {
	cfg := readConfig(configName)
	if cfg == nil {
		cfg = &config{
			Host:      "localhost",
			Port:      "3456",
			UseWs:     true,
			DB:        "boltdb",
			DBFile:    "bolt.db",
			ProofBits: 3,
			ProofMax:  100,
		}
	}

	a.host = cfg.Host
	a.port = cfg.Port

	proof := block.NewProof(cfg.ProofBits, cfg.ProofMax)

	var store block.Storage
	switch cfg.DB {
	case "boltdb":
		store = storage.NewBolt(cfg.DBFile)
	case "buntdb":
		store = storage.NewBuntDB(cfg.DBFile)
	default:
		store = storage.NewBolt(cfg.DBFile)
	}

	a.bc = block.NewBlockchain(proof, store)
	return nil
}

// Run ...
func (a *App) Run() error {
	http.HandleFunc("/", a.home)
	http.HandleFunc("/blocks", a.blocks)
	http.HandleFunc("/mine", a.mineBlock)
	http.HandleFunc("/peers", a.peers)
	http.HandleFunc("/ws", a.ws)

	go a.hub.run()
	return http.ListenAndServe(a.host+":"+a.port, nil)
}

// Stop ...
func (a *App) Stop() error {
	return a.bc.Stop()
}

type config struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	UseWs     bool   `json:"use_ws"`
	DB        string `json:"db"`
	DBFile    string `json:"db_file"`
	ProofBits int    `json:"proof_bits"`
	ProofMax  int    `json:"proof_max"`
}

func readConfig(filename string) *config {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var c config
	err = json.Unmarshal(raw, &c)
	if err != nil && err != io.EOF {
		panic(err)
	}
	return &c
}
