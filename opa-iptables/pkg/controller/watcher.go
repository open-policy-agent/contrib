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

func (w *watcher) watch(workerCh chan<- state) {
	w.mu.RLock()
	for queryPath := range w.watcherState {
		s, ok := w.watcherState[queryPath]
		if !ok {
			w.logger.Debugf("key %v don't exists in watcherState map", queryPath)
			continue
		}
		workerCh <- *s
	}
	w.mu.RUnlock()
}

func (c *Controller) newWatcher() {
	
	workerCh := make(chan state)
	workerDoneCh := make(chan struct{}, c.watcherWorkerCount)
	c.startWorker(workerCh, workerDoneCh)

	c.logger.Info("starting watcher")
	ticker := time.NewTicker(c.w.watcherInterval)
	for {
		select {
		case <-ticker.C:
			c.logger.Debug("checking for any update")
			go c.w.watch(workerCh)
		case <-c.w.watcherDoneCh:
			close(workerCh)
			c.stopWorker(workerDoneCh)
			close(c.w.watcherDoneCh)
			return
		}
	}
}

func (c *Controller) stopWatcher(ctx context.Context) error {
	c.w.watcherDoneCh <- struct{}{}
	for {
		select {
		case <-ctx.Done():
			return context.DeadlineExceeded
		case <-c.w.watcherDoneCh:
			return nil
		}
	}
}

func (c *Controller) startWorker(workerCh <-chan state, done chan<- struct{}) {
	for i := 1 ; i <= c.watcherWorkerCount ; i++ {
		go c.worker(i, workerCh, done)
	}
}

func (c *Controller) stopWorker(done chan struct{}) {
	for i := 1 ; i <= c.watcherWorkerCount ; i++ {
		<-done
	}
	close(done)
}

// worker runs in it's own goroutine.
func (c *Controller) worker(id int, workerCh <-chan state, done chan<- struct{}) {
	c.logger.Infof("Worker %v started", id)

	for s := range workerCh {

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
	c.logger.Infof("worker %v stopped", id)
	done <- struct{}{}
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
