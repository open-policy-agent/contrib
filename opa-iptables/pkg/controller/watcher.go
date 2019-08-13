package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/open-policy-agent/contrib/opa-iptables/pkg/iptables"
)

func (w *watcher) addState(s *state) {
	w.mu.Lock()
	w.watcherState[s.queryPath] = s
	w.mu.Unlock()
}

func (w *watcher) removeState(queryPath string) {
	w.mu.Lock()
	_, ok := w.watcherState[queryPath]
	if ok {
		delete(w.watcherState, queryPath)
	}
	w.mu.Unlock()
}

func (w *watcher) getState(key string) (state, error) {
	w.mu.RLock()
	s, ok := w.watcherState[key]
	if !ok {
		w.mu.RUnlock()
		return state{}, fmt.Errorf("no state found of key %v", key)
	}
	w.mu.RUnlock()
	return *s, nil
}

func (w *watcher) watcher(watchCh chan<- string) {
	w.mu.RLock()
	for k := range w.watcherState {
		watchCh <- k
	}
	w.mu.RUnlock()
}

func (c *Controller) watch() {

	watcherCh := make(chan string)
	for i := 1; i <= 3; i++ {
		go c.worker(i, watcherCh)
	}

	c.logger.Info("starting watcher")
	ticker := time.NewTicker(c.w.watcherInterval)
	for {
		select {
		case <-ticker.C:
			c.logger.Debug("checking for any update")
			go c.w.watcher(watcherCh)
		case <-c.w.done:
			close(watcherCh)
			close(c.w.done)
			return
		}
	}
}

func (c *Controller) shutdownWatcher(ctx context.Context) error {
	c.w.done <- struct{}{}
	for {
		select {
		case <-ctx.Done():
			return context.DeadlineExceeded
		case <-c.w.done:
			return nil
		}
	}
}

// worker runs in it's own goroutine.
func (c *Controller) worker(id int, watchCh <-chan string) {
	c.logger.Infof("Worker %v started", id)

	for queryPath := range watchCh {

		c.w.mu.RLock()
		s, ok := c.w.watcherState[queryPath]
		if !ok {
			c.w.mu.RUnlock()
			c.logger.Debugf("[Worker: %v] key %v don't exists in watcherState map", id, queryPath)
			continue
		}
		c.w.mu.RUnlock()

		res, err := c.handleQuery(s.queryPath, s.payload.Input)
		if err != nil {
			c.logger.Debugf("[Worker: %v] Error while querying opa: %v", id, err)
			continue
		}

		ruleSets, err := iptables.UnmarshalRuleset(res)
		if err != nil {
			c.logger.Debugf("[Worker: %v] Error while Unmarshaling ruleset: %v", id, err)
			continue
		}

		if len(ruleSets) == 1 {
			ruleset := ruleSets[0]
			newID := ruleset.Metadata.ID
			currentID := s.id

			if currentID != newID {

				c.logger.Infof("[Worker: %v] Data changes of queryPath %v, Replacing rules", id, s.queryPath)
				oldRules, _ := c.getCurrentRulesFromOPA(currentID)
				newRules := ruleset.Rules
				err := replaceRules(oldRules, newRules)
				if err != nil {
					c.logger.Error(err)
					continue
				}

				c.putNewRulesToOPA(newID, newRules)
				c.deleteOldRulesFromOPA(currentID)

				newState := state{
					id:        newID,
					payload:   s.payload,
					queryPath: s.queryPath,
				}
				c.w.addState(&newState)
			}
		}
	}
	c.logger.Debugf("worker %v stopped", id)
}

func replaceRules(old, new []iptables.Rule) error {
	//deletes old rules
	err := deleteRules(old)
	if err != nil {
		return err
	}
	//inserts new rules
	err = insertRules(new)
	if err != nil {
		return err
	}
	return nil
}
