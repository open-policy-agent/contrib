package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/open-policy-agent/contrib/opa-iptables/pkg/converter"
	"github.com/open-policy-agent/contrib/opa-iptables/pkg/iptables"
)

func (c *Controller) jsonRuleHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()

		jsonRules, err := converter.IPTableToJSON(r.Body)
		if err != nil {
			c.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var buf bytes.Buffer
		for i, rule := range jsonRules {
			// for first rule
			if i == 0 {
				buf.Write([]byte("["))
			}

			// add each rule to buffer
			buf.WriteString(rule)

			// add ',' and '\n' to each rule except last rule
			if i != len(jsonRules)-1 {
				buf.WriteByte(',')
				buf.WriteByte('\n')
			}

			// for last rule
			if i == len(jsonRules)-1 {
				buf.Write([]byte("]"))
			}
		}
		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
	}
}

func (c *Controller) insertHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ruleSet, err := c.handlePayload(r)
		if err != nil {
			c.logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(ruleSet.Rules) > 0 {
			c.insertRules(ruleSet)
		} else {
			c.logger.Error("Query didn't returned any IPTables rules")
		}
	}
}

func (c *Controller) deleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ruleSet, err := c.handlePayload(r)
		if err != nil {
			c.logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(ruleSet.Rules) > 0 {
			c.deleteRules(ruleSet)
		} else {
			c.logger.Error("Query didn't returned any IPTables rules")
		}
	}
}

type Payload struct {
	Input interface{} `json:"input"`
}

func (c *Controller) handlePayload(r *http.Request) (iptables.RuleSet, error) {
	c.logger.Infof("msg=\"Received Request\" req_method=%v req_path=%v\n", r.Method, r.URL)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return iptables.RuleSet{}, fmt.Errorf("Error while reading body: %v", err)
	}
	defer r.Body.Close()

	var payload Payload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return iptables.RuleSet{}, fmt.Errorf("Error while unmarshalling webhook payload :%v", err)
	}

	queryPath := strings.TrimPrefix(r.FormValue("q"), "/")
	res, err := c.handleQuery(queryPath, payload.Input)
	if err != nil {
		return iptables.RuleSet{}, fmt.Errorf("Error while quering OPA: %v", err)
	}

	if len(string(res)) == 2 && string(res) == "{}" {
		return iptables.RuleSet{}, fmt.Errorf("Provided query path \"%v\" is invalid or not exists", queryPath)
	}

	ruleSet, err := iptables.UnmarshalRules(res)
	if err != nil {
		return iptables.RuleSet{}, fmt.Errorf("Error while Unmarshaling rules: %v", err)
	}
	return ruleSet, nil
}

func (c *Controller) listRules() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		table := mux.Vars(r)["table"]
		if table == "" {
			table = "filter"
		}
		chain := mux.Vars(r)["chain"]
		if chain == "" {
			chain = "INPUT"
		}
		rules, err := iptables.ListRules(table, chain)
		if err != nil {
			c.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for _, rule := range rules {
			fmt.Fprintln(w, rule)
		}
	}
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
