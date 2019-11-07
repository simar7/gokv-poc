package benchmarks

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	bolt "go.etcd.io/bbolt"
)

const NumGoRoutines = 1000

func BoltBatchGoRoutines() {
	f, err := ioutil.TempFile("", "boltdb-BoltBatchGoRoutines-*")
	if err != nil {
		log.Fatal(err)
	}

	// Open the database.
	db, err := bolt.Open(f.Name(), 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(db.Path())

	var wg sync.WaitGroup
	wg.Add(NumGoRoutines)

	// Start a multiple write transaction.
	for i := 0; i < NumGoRoutines; i++ {
		go func(grn int) {
			defer wg.Done()
			if err := db.Batch(func(tx *bolt.Tx) error {
				// Create a bucket.
				b, err := tx.CreateBucket([]byte(fmt.Sprintf("widgets%d", grn)))
				if err != nil {
					return err
				}

				// Set the value "bar" for the key "foo".
				if err := b.Put([]byte(fmt.Sprintf("foo%d", grn)), []byte("bar")); err != nil {
					return err
				}
				return nil
			}); err != nil {
				log.Fatal(err)
			}
		}(i)
	}

	wg.Wait()
	// Read value back in a different read-only transaction.
	//if err := db.View(func(tx *bolt.Tx) error {
	//	value := tx.Bucket([]byte("widgets")).Get([]byte("foo"))
	//	fmt.Printf("[batch] The value of 'foo' is: %s\n", value)
	//	return nil
	//}); err != nil {
	//	log.Fatal(err)
	//}

	// Close database to release file lock.
	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
}

func BoltUpdateGoRoutines() {
	f, err := ioutil.TempFile("", "boltdb-BoltUpdateGoRoutines-*")
	if err != nil {
		log.Fatal(err)
	}

	// Open the database.
	db, err := bolt.Open(f.Name(), 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(db.Path())

	var wg sync.WaitGroup
	wg.Add(NumGoRoutines)

	// Start a multiple write transaction.
	for i := 0; i < NumGoRoutines; i++ {
		go func(grn int) {
			defer wg.Done()
			if err := db.Update(func(tx *bolt.Tx) error {
				// Create a bucket.
				b, err := tx.CreateBucket([]byte(fmt.Sprintf("widgets%d", grn)))
				if err != nil {
					return err
				}

				// Set the value "bar" for the key "foo".
				if err := b.Put([]byte(fmt.Sprintf("foo%d", grn)), []byte("bar")); err != nil {
					return err
				}
				return nil
			}); err != nil {
				log.Fatal(err)
			}
		}(i)
	}

	wg.Wait()
	// Read value back in a different read-only transaction.
	//if err := db.View(func(tx *bolt.Tx) error {
	//	value := tx.Bucket([]byte("widgets")).Get([]byte("foo"))
	//	fmt.Printf("[batch] The value of 'foo' is: %s\n", value)
	//	return nil
	//}); err != nil {
	//	log.Fatal(err)
	//}

	// Close database to release file lock.
	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
}

func BoltBatch() {
	f, err := ioutil.TempFile("", "boltdb-*")
	if err != nil {
		log.Fatal(err)
	}

	// Open the database.
	db, err := bolt.Open(f.Name(), 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(db.Path())

	// Start a write transaction.
	if err := db.Batch(func(tx *bolt.Tx) error {
		// Create a bucket.
		b, err := tx.CreateBucket([]byte("widgets"))
		if err != nil {
			return err
		}

		// Set the value "bar" for the key "foo".
		if err := b.Put([]byte("foo"), []byte("bar")); err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	// Read value back in a different read-only transaction.
	//if err := db.View(func(tx *bolt.Tx) error {
	//	value := tx.Bucket([]byte("widgets")).Get([]byte("foo"))
	//	fmt.Printf("[batch] The value of 'foo' is: %s\n", value)
	//	return nil
	//}); err != nil {
	//	log.Fatal(err)
	//}

	// Close database to release file lock.
	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
}

func BoltUpdate() {
	f, err := ioutil.TempFile("", "boltdb-*")
	if err != nil {
		log.Fatal(err)
	}

	// Open the database.
	db, err := bolt.Open(f.Name(), 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(db.Path())

	// Start a write transaction.
	if err := db.Update(func(tx *bolt.Tx) error {
		// Create a bucket.
		b, err := tx.CreateBucket([]byte("widgets"))
		if err != nil {
			return err
		}

		// Set the value "bar" for the key "foo".
		if err := b.Put([]byte("foo"), []byte("bar")); err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	// Read value back in a different read-only transaction.
	//if err := db.View(func(tx *bolt.Tx) error {
	//	value := tx.Bucket([]byte("widgets")).Get([]byte("foo"))
	//	fmt.Printf("[update] The value of 'foo' is: %s\n", value)
	//	return nil
	//}); err != nil {
	//	log.Fatal(err)
	//}

	// Close database to release file lock.
	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
}
