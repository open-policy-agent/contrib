package iptables

import (
	"reflect"
	"encoding/json"
)

type RuleSet struct {
	Metadata struct {
		ID string `json:"_id"`
	} `json:"metadata"`
	Rules []Rule `json:"rules"`
}

type OpaResponse struct {
	RuleSets []RuleSet `json:"result"`
}

func (or OpaResponse) isEmpty() bool {
	for _, ruleSet := range or.RuleSets {
		if !reflect.DeepEqual(ruleSet,RuleSet{}) {
			return false
		}
	}
	return true
}

func UnmarshalRuleset(opaQueryRes []byte) ([]RuleSet, error) {
	var or OpaResponse
	err := json.Unmarshal(opaQueryRes, &or)
	if err != nil {
		return nil, err
	}
	if or.isEmpty() {
		return []RuleSet{},nil
	}
	return or.RuleSets, nil
}