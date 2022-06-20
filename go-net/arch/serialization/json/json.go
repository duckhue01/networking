package json

import (
	"encoding/json"
	"log"
	"os"

	seri "github.com/duckhue01/lang/go/net/arch/serialization"
)

const (
	jsonFile = "./data/chores.json"
)

type JSON struct {
}

func (r *JSON) Decode() (chores []*seri.Chore, err error) {
	if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
		return chores, err
	}
	df, err := os.Open(jsonFile)
	if err != nil {
		return chores, err
	}
	defer func() {
		if err := df.Close(); err != nil {
			log.Fatalf("closing data file: %v", err)
		}
	}()

	return chores, json.NewDecoder(df).Decode(&chores)
}

func (r *JSON) Encode([]*seri.Chore) (err error) {

	return nil
}
