package database

import (
	"archive/zip"
	"bytes"
	"fmt"
	"path/filepath"

	bolt "go.etcd.io/bbolt"
	"gopkg.in/yaml.v2"
)

func initBuckets(tx *bolt.Tx) error {
	_, err := tx.CreateBucketIfNotExists([]byte("Problems"))
	return err
}

type jutgeDB struct {
	dbFile string
	db     *bolt.DB
	ro     bool
}

// NewJutgeDB returns jutgeDB object
func NewJutgeDB(dbFile string) *jutgeDB {
	return &jutgeDB{dbFile, nil, true}
}

type problemData struct {
	ID           string `yaml:"problem_id"`
	Submission   string `yaml:"submission_id"`
	Compiler     string `yaml:"compiler_id"`
	Annotation   string `yaml:"annotation"`
	Veredict     string `yaml:"veredict"`
	VeredictInfo string `yaml:"veredict_info"`
	Title        string `yaml:"title"`
}

func (jDB *jutgeDB) initRead() (err error) {
	if jDB.db == nil {
		jDB.db, err = bolt.Open(jDB.dbFile, 0600, &bolt.Options{ReadOnly: true})
		jDB.ro = true
	}
	return
}

func (jDB *jutgeDB) initWrite() (err error) {
	if jDB.db == nil || jDB.ro {
        if jDB != nil {
            jDB.db.Close()
        }
		jDB.db, err = bolt.Open(jDB.dbFile, 0600, nil)
		jDB.ro = false
		jDB.db.Update(initBuckets)
	}
	return
}

func (jDB *jutgeDB) Close() {
	if jDB.db != nil {
		jDB.db.Close()
	}
}

func (jDB *jutgeDB) Print() {
	jDB.initRead()
	jDB.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Problems"))

		b.ForEach(func(k, v []byte) error {
			fmt.Printf("code=%s, title=\"%s\"\n", k, v)
			return nil
		})

		return nil
	})
}

func (jDB *jutgeDB) Query(code string) (string, error) {
	jDB.initRead()
	title := ""
	err := jDB.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Problems"))
		v := b.Get([]byte(code))
		title = string(v)
		return nil
	})

	return title, err
}

func (jDB *jutgeDB) Add(code, title string) error {
	jDB.initWrite()
	return jDB.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Problems"))
		return b.Put([]byte(code), []byte(title))
	})
}

func (jDB *jutgeDB) ImportZip(filename string) error {
	r, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	defer r.Close()

	jDB.initWrite()

	for _, f := range r.File {
		if filepath.Ext(f.Name) == ".txt" {
			rc, err := f.Open()
			if err != nil {
				return err
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(rc)

			data := problemData{}

			err = yaml.Unmarshal(buf.Bytes(), &data)
			if err != nil {
				if filepath.Base(f.Name) == "annotation.txt" {
					continue
				}
				return err
			}

			// TODO: do in batch, this is slow af
			err = jDB.db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("Problems"))
				return b.Put([]byte(data.ID), []byte(data.Title))
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
