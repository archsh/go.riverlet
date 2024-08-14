package riverlet

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/glebarez/go-sqlite"
	"strings"
	"time"
)

const (
	tableName   = "jobs"
	tableSchema = `
CREATE TABLE IF NOT EXISTS ` + tableName + ` (
	ID INTEGER PRIMARY KEY   AUTOINCREMENT,
	Identity CHAR(256) NOT NULL,
	Args TEXT NOT NULL,
	Status INTEGER NOT NULL DEFAULT 0,
	Attempts INTEGER NOT NULL DEFAULT 0,
	MaxAttempts INTEGER NOT NULL DEFAULT 1,
	Priority INTEGER NOT NULL DEFAULT 1,
	Created INTEGER NULL,
	Updated INTEGER NULL
)`
)

type Persistence[T Argument] interface {
	Insert(job T, priority int, maxAttempts int) (*JobObject[T], error)
	Delete(seqs ...int64) (int64, error)
	Get(seq int64) (*JobObject[T], error)
	Select(queries map[string]interface{}, offset int, limit int, orderBy ...string) ([]*JobObject[T], int64, error)
	Update(seq int64, vals map[string]interface{}) (int64, error)
}

type inMemoryPersistence[T Argument] struct {
	db  *sql.DB
	seq int64
}

func (p *inMemoryPersistence[T]) SeqNext() int64 {
	p.seq++
	return p.seq
}

func (p *inMemoryPersistence[T]) Insert(arg T, priority int, maxAttempts int) (*JobObject[T], error) {
	var j JobObject[T]
	var t = time.Now()
	j.Seq = p.SeqNext()
	j.Identity = arg.Identity()
	j.Args = arg
	j.Status = AVAILABLE
	j.MaxAttempts = maxAttempts
	j.Priority = priority
	j.Created = &t
	if bs, ee := json.Marshal(j.Args); nil != ee {
		return nil, ee
	} else if tx, e := p.db.Begin(); nil != e {
		return nil, e
	} else if r, e := tx.Exec("INSERT INTO "+tableName+"(Identity,Args,Status,Priority,MaxAttempts,Created) VALUES($1,$2,$3,$4,$5,$6)",
		j.Identity, string(bs), j.Status, j.Priority, j.MaxAttempts, j.Created.UnixMilli()); nil != e {
		return nil, e
	} else if n, e := r.LastInsertId(); nil != e {
		return nil, e
	} else if e := tx.Commit(); nil != e {
		return nil, e
	} else {
		j.Seq = n
	}
	return &j, nil
}

func (p *inMemoryPersistence[T]) Delete(seqs ...int64) (int64, error) {
	if len(seqs) < 1 {
		return 0, nil
	}
	var ids []string = make([]string, len(seqs))
	for i, n := range seqs {
		ids[i] = fmt.Sprint(n)
	}

	if tx, e := p.db.Begin(); nil != e {
		return 0, e
	} else if r, e := tx.Exec("DELETE FROM " + tableName + "WHERE ID IN (" + strings.Join(ids, ",") + ")"); nil != e {
		return 0, e
	} else if e := tx.Commit(); nil != e {
		return 0, e
	} else {
		return r.RowsAffected()
	}
}

func (p *inMemoryPersistence[T]) Update(seq int64, vals map[string]interface{}) (int64, error) {
	return 0, nil
}

func (p *inMemoryPersistence[T]) scan(row *sql.Row, j *JobObject[T]) error {
	var args string
	var created, updated sql.NullInt64
	if e := row.Scan(&j.Seq, &j.Identity, &args, &j.Status, &j.Attempts, &j.MaxAttempts, &j.Priority, &created, &updated); nil != e {
		return e
	} else if e := json.Unmarshal([]byte(args), &j.Args); nil != e {
		return e
	}
	if created.Valid {
		var t = time.UnixMilli(created.Int64)
		j.Created = &t
	}
	if updated.Valid {
		var t = time.UnixMilli(updated.Int64)
		j.Updated = &t
	}
	return nil
}

func (p *inMemoryPersistence[T]) Get(seq int64) (*JobObject[T], error) {
	row := p.db.QueryRow("SELECT ID,Identity,Args,Status,Attempts,MaxAttempts,Priority,Created,Updated FROM "+tableName+" WHERE ID = $1", seq)
	var j JobObject[T]

	if nil == row {
		return nil, fmt.Errorf("not found ID: %d", seq)
	} else if nil != row.Err() {
		return nil, row.Err()
	} else if e := p.scan(row, &j); nil != e {
		return nil, e
	}
	return &j, nil
}

func (p *inMemoryPersistence[T]) Select(queries map[string]interface{}, offset int, limit int, orderBy ...string) ([]*JobObject[T], int64, error) {
	return nil, 0, nil
}

func NewBuiltinPersistence[T Argument]() Persistence[T] {
	// Create a new database
	db, err := sql.Open("sqlite", ":memory:") //
	if err != nil {
		panic(err)
	}
	if tx, e := db.Begin(); nil != e {
		panic(e)
	} else if _, e := tx.Exec(tableSchema); nil != e {
		panic(e)
	} else if e := tx.Commit(); nil != e {
		panic(e)
	}
	return &inMemoryPersistence[T]{db: db}
}
