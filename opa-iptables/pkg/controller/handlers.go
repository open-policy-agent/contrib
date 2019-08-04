package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gorilla/mux"

	"github.com/open-policy-agent/contrib/opa-iptables/pkg/converter"
	"github.com/open-policy-agent/contrib/opa-iptables/pkg/iptables"
)

func (c *Controller) jsonRuleHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.logger.Infof("msg=\"Received Request\" req_method=%v req_path=%v\n", r.Method, r.URL)
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

// insertRuleHandler query OPA using provided payload through request and get iptables rules
// and insert them to the kernel
//
//      Server Response:
//
//      200 OK           - 	 Successfully inserted given iptables rules
//      400 Bad Request  -   If provided query path didn't resolve to any defined OPA policy
//                           rule or server fail to parse JSON payload
//      404 Not Found    -   OPA policy rule didn't return any iptables rules
//      500 Server Error -   Fail to insert given iptables rules
//
//
func (c *Controller) insertRuleHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ruleSet, err := c.handlePayload(r)
		if err != nil {
			c.logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(ruleSet.Rules) > 0 {
			err := c.insertRules(ruleSet)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "%s. Checkout log for more details", err)
			}
		} else {
			c.logger.Error("Query didn't returned any IPTables rules")
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

// deleteRuleHandler query OPA using provided payload through request and get iptables rules
// and delete them from the kernel
//
//      Server Response:
//
//      200 OK           - 	 Successfully deleted given iptables rules
//      400 Bad Request  -   If provided query path didn't resolve to any defined OPA policy
//                           rule or server fail to parse JSON payload
//      404 Not Found    -   OPA policy rule didn't return any iptables rules
//      500 Server Error -   Fail to delete given iptables rules
//
//
func (c *Controller) deleteRuleHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ruleSet, err := c.handlePayload(r)
		if err != nil {
			c.logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(ruleSet.Rules) > 0 {
			err := c.deleteRules(ruleSet)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "%s. Checkout log for more details", err)
			}
		} else {
			c.logger.Error("Query didn't returned any IPTables rules")
			w.WriteHeader(http.StatusNotFound)
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
		return iptables.RuleSet{}, fmt.Errorf("Provided query path \"%v\" is not valid path to policy rule", queryPath)
	}

	ruleSet, err := iptables.UnmarshalRules(res)
	if err != nil {
		return iptables.RuleSet{}, fmt.Errorf("Error while Unmarshaling rules: %v", err)
	}
	return ruleSet, nil
}

func (c *Controller) listRulesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.logger.Infof("msg=\"Received Request\" req_method=%v req_path=%v\n", r.Method, r.URL)
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

func (c *Controller) listAllRulesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.logger.Infof("msg=\"Received Request\" req_method=%v req_path=%v\n", r.Method, r.URL)
		verboseStr := r.FormValue("verbose")
		verbose := stringToBool(verboseStr)
		var iptableTableList = [...]string{"filter", "nat"}

		if verbose {
			var buf bytes.Buffer
			for _, table := range iptableTableList {
				stdout, err := runCommand("/sbin/iptables", "-n", "-v", "-L", "-t", table)
				if err != nil {
					c.logger.Errorf("Unable to list iptables rule: %v",err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				buf.Write(stdout)
			}
			fmt.Fprint(w, buf.String())
		} else {
			stdout, err := runCommand("/sbin/iptables", "-S")
			if err != nil {
				c.logger.Errorf("Unable to list iptables rule: %v",err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			fmt.Fprint(w, string(stdout))
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

func stringToBool(value string) bool {
	if value == "true" {
		return true
	}
	return false
}

func runCommand(name string, args ...string) (output []byte, err error) {
	cmd := exec.Command(name, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	stderrorCmd, err := ioutil.ReadAll(stderr)
	if err != nil {
		return nil, err
	}

	stdoutCmd, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	if string(stderrorCmd) != "" {
		return nil, fmt.Errorf(string(stderrorCmd))
	}

	return stdoutCmd, nil
}