package common

import (
	"gopkg.in/fsnotify.v1"
	"os"
	"path/filepath"
)

type Watcher struct {
	Events <-chan fsnotify.Event
	Errors <-chan error

	sendEvents  chan<- fsnotify.Event
	sendErrors  chan<- error
	fswatcher   *fsnotify.Watcher
	watchingDir map[string]interface{}
}

func NewWatcher(dir string) (*Watcher, error) {
	fswatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	ret := &Watcher{
		fswatcher:   fswatcher,
		watchingDir: make(map[string]interface{}),
	}

	chEvents := make(chan fsnotify.Event)
	chErrors := make(chan error)
	ret.Events = chEvents
	ret.sendEvents = chEvents
	ret.Errors = chErrors
	ret.sendErrors = chErrors

	go func() {

		for {
			select {
			case ev, ok := <-fswatcher.Events:
				if !ok {
					break
				}
				_, watching := ret.watchingDir[ev.Name]
				var err error
				st, statErr := os.Stat(ev.Name)
				if statErr == nil && st.IsDir() {
					if !watching {
						err = ret.Add(ev.Name)
					}
				} else if watching {
					err = ret.Remove(ev.Name)
				}

				if err != nil {
					chErrors <- err
				}

				chEvents <- ev

			case err, ok := <-fswatcher.Errors:
				if !ok {
					break
				}
				chErrors <- err
			}
		}
	}()

	err = filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				if e := ret.Add(path); e != nil {
					return e
				}
			}

			return nil
		})

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (self *Watcher) Add(path string) (err error) {
	if _, ok := self.watchingDir[path]; !ok {
		self.watchingDir[path] = true
		err = self.fswatcher.Add(path)
	}
	return
}

func (self *Watcher) Remove(path string) (err error) {
	if _, ok := self.watchingDir[path]; ok {
		delete(self.watchingDir, path)
		if _, statErr := os.Stat(path); statErr == nil {
			err = self.fswatcher.Remove(path)
		}
	}
	return
}

func (self *Watcher) Close() {
	self.fswatcher.Close()
	close(self.sendEvents)
	close(self.sendErrors)
}
