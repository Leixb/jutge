package database

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/imroc/req"
	bolt "go.etcd.io/bbolt"
	"gopkg.in/yaml.v2"
)

type jutgeDB struct {
	dbFile string
	db     *bolt.DB
	ro     bool
}

// NewJutgeDB returns jutgeDB object
func NewJutgeDB(dbFile string) *jutgeDB {
	jDB := jutgeDB{dbFile, nil, true}
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(dbFile), 0755)
		if err != nil {
			panic(err)
		}
		_, err = os.Create(dbFile)
		if err != nil {
			panic(err)
		}

		full_path, _ := filepath.Abs(dbFile)
		fmt.Println("Created new DB file", full_path)
		jDB.initWrite()
	}
	return &jDB
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
		if jDB.db != nil {
			jDB.db.Close()
		}
		jDB.db, err = bolt.Open(jDB.dbFile, 0600, nil)
		if err != nil {
			return err
		}

		jDB.ro = false
		jDB.db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte("Problems"))
			return err
		})
	}
	return
}

func (jDB *jutgeDB) Close() {
	if jDB.db != nil {
		jDB.db.Close()
	}
}

func (jDB *jutgeDB) Print() error {
	err := jDB.initRead()
	if err != nil {
		return err
	}

	return jDB.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Problems"))

		b.ForEach(func(k, v []byte) error {
			fmt.Printf("code=%s, title=\"%s\"\n", k, v)
			return nil
		})

		return nil
	})
}

func (jDB *jutgeDB) Query(code string) (string, error) {
	err := jDB.initRead()
	if err != nil {
		return "", err
	}

	title := ""
	err = jDB.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Problems"))
		v := b.Get([]byte(code))
		title = string(v)
		return nil
	})

	return title, err
}

func (jDB *jutgeDB) Add(code, title string) error {
	err := jDB.initWrite()
	if err != nil {
		return err
	}

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

func (jDB *jutgeDB) Download() error {
	if jDB.db != nil {
		jDB.db.Close()
	}
	url := "https://raw.githubusercontent.com/Leixb/jutge/master/jutge.db"
	r, err := req.Get(url)
	if err != nil {
		return err
	}

	return r.ToFile(jDB.dbFile)
}
