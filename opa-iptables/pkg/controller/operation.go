package controller

import (
	"github.com/open-policy-agent/contrib/opa-iptables/pkg/iptables"
)

func (c *Controller) insertRules(ruleSet iptables.RuleSet) {
	successCount := 0
	totalRules := len(ruleSet.Rules)
	for _, rule := range ruleSet.Rules {
		c.logger.Debugf("Inserting Rule: %v", rule.String())
		err := rule.AddRule()
		if err != nil {
			c.logger.Errorf("Error while inserting above rule: %v", err)
			continue
		}
		successCount = successCount + 1
	}
	c.logger.Infof("Inserted %v out of %v rules (%v/%v)", successCount, totalRules,successCount,totalRules)
}

func (c *Controller) deleteRules(ruleSet iptables.RuleSet) {
	successCount := 0
	totalRules := len(ruleSet.Rules)
	for _, rule := range ruleSet.Rules {
		c.logger.Debugf("Deleting Rule: %v", rule.String())
		err := rule.DeleteRule()
		if err != nil {
			c.logger.Errorf("Error while deleting rule: %v", err)
			continue
		}
		successCount = successCount + 1
	}
	c.logger.Infof("Deleted %v out of %v rules (%v/%v)", successCount, totalRules, successCount, totalRules)
}

func (c *Controller) testRules(ruleSet iptables.RuleSet) {
	for i, rule := range ruleSet.Rules {
		c.logger.Infof("Rule %v: %v\n", i+1, rule.String())
	}
}
