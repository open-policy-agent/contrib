package controller

import (
	"encoding/json"

	"github.com/open-policy-agent/contrib/opa-iptables/pkg/iptables"
)

func (c *Controller) putNewRulesToOPA(id string, rules []iptables.Rule) error {
	data, err := iptables.MarshalRules(rules)
	if err != nil {
		return err
	}
	return c.opaClient.PutData("state/"+id, data)
}

func (c *Controller) deleteOldRulesFromOPA(id string) error {
	return c.opaClient.DeleteData("state/"+id)
}

func (c *Controller) getCurrentRulesFromOPA(id string) ([]iptables.Rule,error) {
	data, err := c.opaClient.GetData("state/" + id)
	if err != nil {
		return nil, err
	}

	rules, _ := iptables.UnmarshalRules(data)
	return rules, nil
}

func (c *Controller) handleQuery(path string, data interface{}) ([]byte, error) {
	input, err := marshalInput(data)
	if err != nil {
		return nil, err
	}

	res, err := c.opaClient.DoQuery(path, input)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func marshalInput(data interface{}) ([]byte, error) {
	inputMap := make(map[string]interface{})
	inputMap["input"] = data
	return json.Marshal(inputMap)
}
