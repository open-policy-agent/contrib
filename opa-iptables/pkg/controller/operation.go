package controller

import (
	"github.com/open-policy-agent/contrib/opa-iptables/pkg/logging"
	"fmt"
	"github.com/open-policy-agent/contrib/opa-iptables/pkg/iptables"
)

func insertRules(rules []iptables.Rule) error {
	logger := logging.GetLogger()
	successCount := 0
	totalRules := len(rules)
	var gotError bool

	for _, rule := range rules {
		logger.Debugf("Inserting Rule: %v", rule.String())
		err := rule.AddRule()
		if err != nil {
			gotError = true
			logger.Errorf("Error while inserting rule: %v", err)
			continue
		}
		successCount++
	}
	
	logger.Infof("Inserted %v out of %v rules (%v/%v)", successCount, totalRules,successCount,totalRules)
	if gotError {
		return fmt.Errorf("get error during inserting rules")
	}
	return nil
}

func deleteRules(rules []iptables.Rule) error {
	logger := logging.GetLogger()
	successCount := 0
	totalRules := len(rules)
	var gotError bool

	for _, rule := range rules {
		logger.Debugf("Deleting Rule: %v", rule.String())
		err := rule.DeleteRule()
		if err != nil {
			gotError = true
			logger.Errorf("Error while deleting rule: %v", err)
			continue
		}
		successCount++
	}

	logger.Infof("Deleted %v out of %v rules (%v/%v)", successCount, totalRules, successCount, totalRules)
	if gotError {
		return fmt.Errorf("get error during deleting rules")
	}
	return nil
}

func testRules(ruleSet iptables.RuleSet) {
	logger := logging.GetLogger()
	for i, rule := range ruleSet.Rules {
		logger.Infof("Rule %v: %v\n", i+1, rule.String())
	}
}