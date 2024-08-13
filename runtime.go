package riverlet

import "database/sql"

type Runtime struct {
	works map[string]Worker[Argument]
	db    *sql.DB
}

func NewRuntime() (*Runtime, error) {
	var r = &Runtime{}
	return r, nil
}

func (r *Runtime) RegisterWorker(w Worker[Argument]) error {
	if nil == r.works {
		r.works = make(map[string]Worker[Argument])
	}
	if _, b := r.works[w.DefArgs().Identity()]; b {
		panic("worker " + w.DefArgs().Identity() + " already registered")
	}
	r.works[w.DefArgs().Identity()] = w
	return nil
}

func (r *Runtime) AddJob(job Job[Argument]) error {
	return nil
}

func (r *Runtime) Start() error {
	return nil
}

func (r *Runtime) Stop() error {
	return nil
}
