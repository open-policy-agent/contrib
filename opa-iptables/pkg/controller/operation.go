package controller

import (
	"fmt"
	"github.com/open-policy-agent/contrib/opa-iptables/pkg/iptables"
)

func (c *Controller) insertRules(ruleSet iptables.RuleSet) error {
	successCount := 0
	totalRules := len(ruleSet.Rules)
	var gotError bool

	for _, rule := range ruleSet.Rules {
		c.logger.Debugf("Inserting Rule: %v", rule.String())
		err := rule.AddRule()
		if err != nil {
			gotError = true
			c.logger.Errorf("Error while inserting rule: %v", err)
			continue
		}
		successCount++
	}
	c.logger.Infof("Inserted %v out of %v rules (%v/%v)", successCount, totalRules,successCount,totalRules)
	if gotError {
		return fmt.Errorf("get error during inserting rules")
	}
	return nil
}

func (c *Controller) deleteRules(ruleSet iptables.RuleSet) error {
	successCount := 0
	totalRules := len(ruleSet.Rules)
	var gotError bool

	for _, rule := range ruleSet.Rules {
		c.logger.Debugf("Deleting Rule: %v", rule.String())
		err := rule.DeleteRule()
		if err != nil {
			gotError = true
			c.logger.Errorf("Error while deleting rule: %v", err)
			continue
		}
		successCount++
	}

	c.logger.Infof("Deleted %v out of %v rules (%v/%v)", successCount, totalRules, successCount, totalRules)
	if gotError {
		return fmt.Errorf("get error during deleting rules")
	}
	return nil
}

func (c *Controller) testRules(ruleSet iptables.RuleSet) {
	for i, rule := range ruleSet.Rules {
		c.logger.Infof("Rule %v: %v\n", i+1, rule.String())
	}
}