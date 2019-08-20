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
	cmd "github.com/open-policy-agent/contrib/opa-iptables/pkg/command"
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

		ruleSets, request, err := c.handlePayload(r)
		if err != nil {
			c.logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		
		var insertError error
		if len(ruleSets) > 0 {
			for _, ruleSet := range ruleSets {

				if len(ruleSet.Rules) > 0 {
					err := insertRules(ruleSet.Rules)
					if err != nil {
						insertError = err
						continue
					}
				}
			}

			if insertError != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "%s. Checkout log for more details", err)
				return
			}

		} else {
			c.logger.Error("Query didn't returned any ruleSet")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if c.experimental {

			watch := stringToBool(r.FormValue("watch"))

			if len(ruleSets) > 1 && watch {
				c.logger.Error("Unable to watch queryPath. Query returns multiple ruleSet")
				return
			}
			if watch {

				rs := ruleSets[0]
				s := state{
					id:        rs.Metadata.ID,
					payload:   request.p,
					queryPath: request.queryPath,
				}

				if s.id == "" {
					c.logger.Error("Unable to watch current queryPath. RuleSet cotains empty \"_id\" field.")
					return
				}

				err := c.putNewRulesToOPA(s.id,rs.Rules)
				if err != nil {
					return
				}
				
				c.w.addState(&s)
			}
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
		ruleSets, request, err := c.handlePayload(r)
		if err != nil {
			c.logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var deleteError error
		if len(ruleSets) > 0 {
			for _, ruleSet := range ruleSets {
				if len(ruleSet.Rules) > 0 {
					err := deleteRules(ruleSet.Rules)
					if err != nil {
						deleteError = err
						continue
					}
				}
			}

			if deleteError != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "%s. Checkout log for more details", err)
				return
			}

			if c.experimental {
				s, err := c.w.getState(request.queryPath)
				if err != nil {
					c.logger.Error(err)
					return
				}
	
				err = c.deleteOldRulesFromOPA(s.id)
				if err != nil {
					c.logger.Error(err)
					return
				}
				c.w.removeState(s.queryPath)
			}

		} else {
			c.logger.Error("Query didn't returned any RuleSet")
			w.WriteHeader(http.StatusNotFound)
		}
	}
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
		
		verbose := stringToBool(r.FormValue("verbose"))
		var iptableTableList = [...]string{"filter", "nat"}

		if verbose {
			var buf bytes.Buffer
			for _, table := range iptableTableList {
				stdout, err := cmd.RunCommand("/sbin/iptables", "-n", "-v", "-L", "-t", table)
				if err != nil {
					c.logger.Errorf("Unable to list iptables rule: %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				buf.Write(stdout)
			}
			fmt.Fprint(w, buf.String())
		} else {
			stdout, err := cmd.RunCommand("/sbin/iptables", "-S")
			if err != nil {
				c.logger.Errorf("Unable to list iptables rule: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			fmt.Fprint(w, string(stdout))
		}
	}
}

func (c *Controller) handlePayload(r *http.Request) ([]iptables.RuleSet, request, error) {
	c.logger.Infof("msg=\"Received Request\" req_method=%v req_path=%v\n", r.Method, r.URL)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, request{}, fmt.Errorf("Error while reading body: %v", err)
	}
	defer r.Body.Close()

	var payload payload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return nil, request{}, fmt.Errorf("Error while unmarshalling payload :%v", err)
	}

	queryPath := strings.TrimPrefix(r.FormValue("q"), "/")
	res, err := c.handleQuery(queryPath, payload.Input)
	if err != nil {
		return nil, request{}, fmt.Errorf("Error while quering OPA: %v", err)
	}

	if len(string(res)) == 2 && string(res) == "{}" {
		return nil, request{}, fmt.Errorf("Provided query path \"%v\" is not valid path to policy rule", queryPath)
	}

	ruleSets, err := iptables.UnmarshalRuleset(res)
	if err != nil {
		return nil, request{}, fmt.Errorf("Error while Unmarshaling ruleset: %v", err)
	}

	return ruleSets, request{queryPath: queryPath, p: payload}, nil
}

func stringToBool(value string) bool {
	if value == "true" {
		return true
	}
	return false
}