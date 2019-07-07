package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/open-policy-agent/contrib/opa-iptables/pkg/iptables"
)

func (c *Controller) webhookHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.logger.Infof("msg=\"Received Request\" req_method=%v req_path=%v\n", r.Method, r.URL)
		if r.Method == http.MethodPost {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				c.logger.Errorf("Error while reading body: %v", err)
				return
			}
			defer r.Body.Close()

			var payload Payload
			err = json.Unmarshal(body, &payload)
			if err != nil {
				c.logger.Errorf("Error while unmarshalling webhook payload :%v", err)
				return
			}

			queryPath := strings.TrimPrefix(payload.QueryPath, "/")
			res, err := c.handleQuery(queryPath, payload.Input)
			if err != nil {
				c.logger.Errorf("Error while quering OPA: %v", err)
				return
			}

			if len(string(res)) == 2 && string(res) == "{}" {
				c.logger.Errorf("Provided query path \"%v\" is not valid", payload.QueryPath)
				return
			}

			ruleSet, err := iptables.UnmarshalRules(res)
			if err != nil {
				c.logger.Errorf("Error while Unmarshaling rules: %v", err)
				return
			}

			if len(ruleSet.Rules) > 0 {
				switch payload.Op {
				case insertOp:
					c.insertRules(ruleSet)
				case deleteOp:
					c.deleteRules(ruleSet)
				case testOp:
					fallthrough
				default:
					c.testRules(ruleSet)
				}
			} else {
				c.logger.Error("Query didn't returned any IPTables rules")
			}

		} else {
			res := "This endpoint don't support provided method"
			resStatusCode := http.StatusMethodNotAllowed
			w.Header().Add("Allow", http.MethodPost)
			w.WriteHeader(resStatusCode)
			w.Write([]byte(res))
			c.logger.Errorf("msg=\"Sent Response\" req_method=%v req_path=%v res_bytes=%v res_status=%v\n", r.Method, r.URL, len(res), resStatusCode)
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
