package southxchange

import (
	"encoding/json"
	"fmt"
)

type MarketSummary struct {
	Coin     string  `json:"Coin"`
	Base     string  `json:"Base"`
}

func (n *MarketSummary) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&n.Coin, &n.Base}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("wrong number of fields in MarketSummary: %d != %d", g, e)
	}
	return nil
}
