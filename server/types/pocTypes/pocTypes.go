package pocTypes

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
	"strings"
)

type InclusionList struct {
	Transactions []bellatrix.Transaction `ssz-max:"1048576,1073741824" ssz-size:"?,?"`
}

// inclusionListJSON is the spec representation of the struct.
type inclusionListJSON struct {
	Transactions []string `json:"transactions"`
}

// inclusionListYAML is the spec representation of the struct.
type inclusionListYAML struct {
	Transactions []string `yaml:"transactions"`
}

// MarshalJSON implements json.Marshaler.
func (e *InclusionList) MarshalJSON() ([]byte, error) {
	transactions := make([]string, len(e.Transactions))
	for i := range e.Transactions {
		transactions[i] = fmt.Sprintf("%#x", e.Transactions[i])
	}

	return json.Marshal(&inclusionListJSON{
		Transactions: transactions,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *InclusionList) UnmarshalJSON(input []byte) error {
	var data inclusionListJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return e.unpack(&data)
}

//nolint:gocyclo
func (e *InclusionList) unpack(data *inclusionListJSON) error {
	if data.Transactions == nil {
		return errors.New("transactions missing")
	}
	transactions := make([]bellatrix.Transaction, len(data.Transactions))
	for i := range data.Transactions {
		if data.Transactions[i] == "" {
			return errors.New("transaction missing")
		}
		tmp, err := hex.DecodeString(strings.TrimPrefix(data.Transactions[i], "0x"))
		if err != nil {
			return errors.Wrap(err, "invalid value for transaction")
		}
		transactions[i] = bellatrix.Transaction(tmp)
	}
	e.Transactions = transactions

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (e *InclusionList) MarshalYAML() ([]byte, error) {
	transactions := make([]string, len(e.Transactions))
	for i := range e.Transactions {
		transactions[i] = fmt.Sprintf("%#x", e.Transactions[i])
	}
	yamlBytes, err := yaml.MarshalWithOptions(&inclusionListYAML{
		Transactions: transactions,
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (e *InclusionList) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var data inclusionListJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return err
	}

	return e.unpack(&data)
}

// String returns a string version of the structure.
func (s *InclusionList) String() string {
	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
