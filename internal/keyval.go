package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"time"
)

type item struct {
	Value      string
	LastAccess time.Time
}

type KeyVal struct {
	fileName string
	ttl      time.Duration
	mux      sync.RWMutex

	values   map[string]*item
	modified chan bool
}

func NewKeyVal(fileName string, ttl time.Duration) *KeyVal {
	kv := &KeyVal{
		fileName: fileName,
		ttl:      ttl,
		modified: make(chan bool),
	}

	// Load values from file
	file, err := os.Open(fileName)
	if err == nil {
		defer file.Close()

		data, _ := ioutil.ReadAll(file)
		_ = json.Unmarshal(data, &kv.values)
	} else {
		kv.values = map[string]*item{}
	}

	go func(kv *KeyVal) {
		for now := range time.Tick(time.Second) {
			didDelete := false

			kv.mux.Lock()
			for key, it := range kv.values {
				if now.Sub(it.LastAccess) > kv.ttl {
					delete(kv.values, key)
					didDelete = true
				}
			}
			kv.mux.Unlock()

			kv.setModified(didDelete)
		}
	}(kv)

	go func(kv *KeyVal) {
		for {
			<-kv.modified

			err := kv.writeValuesToFile()
			if err != nil {
				fmt.Printf("cannot write values to file: %v\n", err)
			}
		}
	}(kv)

	return kv
}

func (kv *KeyVal) Get(key string) (string, bool) {
	kv.mux.RLock()
	defer kv.mux.RUnlock()

	var value string

	it, ok := kv.values[key]
	if ok {
		value = it.Value
		it.LastAccess = time.Now()
	}

	return value, ok
}

func (kv *KeyVal) Put(key string, value string) {
	kv.mux.Lock()
	defer kv.mux.Unlock()

	if _, ok := kv.values[key]; !ok {
		kv.values[key] = &item{}
	}

	kv.values[key].Value = value
	kv.values[key].LastAccess = time.Now()
	kv.setModified(true)
}

func (kv *KeyVal) setModified(modified bool) {
	if modified {
		select {
		case kv.modified <- true:
		default:
		}
	}
}

func (kv *KeyVal) writeValuesToFile() error {
	kv.mux.RLock()
	defer kv.mux.RUnlock()

	data, err := json.MarshalIndent(kv.values, "", "  ")
	if err != nil {
		return err
	}

	dirname := path.Dir(kv.fileName)
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		err = os.MkdirAll(dirname, 0755)
		if err != nil {
			return err
		}
	}

	return ioutil.WriteFile(kv.fileName, data, 0644)
}
