package data

import (
	"bytes"
	"io/ioutil"
	"os"
	"sort"
	"testing"

	"github.com/skycoin/skycoin/src/cipher"
	// "github.com/skycoin/skycoin/src/cipher/encoder"
)

//
// helper functions
//

func shouldNotPanic(t *testing.T) {
	if err := recover(); err != nil {
		t.Error("unexpected panic:", err)
	}
}

func testPath(t *testing.T) string {
	fl, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer fl.Close()
	return fl.Name()
}

func testDriveDB(t *testing.T) (db DB, cleanUp func()) {
	dbFile := testPath(t)
	db, err := NewDriveDB(dbFile)
	if err != nil {
		os.Remove(dbFile)
		t.Fatal(err)
	}
	cleanUp = func() {
		db.Close()
		os.Remove(dbFile)
	}
	return
}

// returns RootPack that contains dummy Root field,
// the field can't be used to encode/decode
func getRootPack(seq uint64, content string) (rp RootPack) {
	rp.Seq = seq
	if seq != 0 {
		rp.Prev = cipher.SumSHA256([]byte("any"))
	}
	rp.Root = []byte(content)
	rp.Hash = cipher.SumSHA256(rp.Root)
	return
}

type testObjectKeyValue struct {
	key   cipher.SHA256
	value []byte
}

type testObjectKeyValues []testObjectKeyValue

func (t testObjectKeyValues) Len() int {
	return len(t)
}

func (t testObjectKeyValues) Less(i, j int) bool {
	return bytes.Compare(t[i].key[:], t[j].key[:]) < 0
}

func (t testObjectKeyValues) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func testSortedObjects(input ...string) (to testObjectKeyValues) {
	to = make(testObjectKeyValues, 0, len(input))
	for _, s := range input {
		to = append(to, testObjectKeyValue{
			key:   cipher.SumSHA256([]byte(s)),
			value: []byte(s),
		})
	}
	sort.Sort(to)
	return
}

// testOrderedPublicKeys retursn slice with two generated
// public keys in ascending order
func testOrderedPublicKeys() []cipher.PubKey {
	// add feeds
	pk1, _ := cipher.GenerateKeyPair()
	pk2, _ := cipher.GenerateKeyPair()

	// be sure that keys are not equal
	for pk1 == pk2 {
		pk2, _ = cipher.GenerateKeyPair()
	}

	// oreder
	if bytes.Compare(pk2[:], pk1[:]) < 0 { // if pk2 < pk1
		pk1, pk2 = pk2, pk1 // swap
	}

	return []cipher.PubKey{pk1, pk2}
}

func testComparePublicKeyLists(t *testing.T, a, b []cipher.PubKey) {
	if len(a) != len(b) {
		t.Error("wrong list length")
		return
	}
	for i, ax := range a {
		if bx := b[i]; bx != ax {
			t.Errorf("wrong item %d: want %s, got %s",
				i,
				ax.Hex(), // shortHex(
				bx.Hex()) // shortHex(
		}
	}
}

//
// Tests
//

//
// DB
//

func testDBView(t *testing.T, db DB) {
	t.Skip("(TODO) not implemenmted yet")
}

func TestDB_View(t *testing.T) {
	// View(func(t Tv) error) (err error)

	t.Run("memory", func(t *testing.T) {
		testDBView(t, NewMemoryDB())
	})

	t.Run("drive", func(t *testing.T) {
		db, cleanUp := testDriveDB(t)
		defer cleanUp()
		testDBView(t, db)
	})

}

func testDBUpdate(t *testing.T, db DB) {
	t.Skip("(TODO) not implemented yet")
}

func TestDB_Update(t *testing.T) {
	// Update(func(t Tu) error) (err error)

	t.Run("memory", func(t *testing.T) {
		testDBUpdate(t, NewMemoryDB())
	})

	t.Run("drive", func(t *testing.T) {
		db, cleanUp := testDriveDB(t)
		defer cleanUp()
		testDBUpdate(t, db)
	})

}

func testDBStat(t *testing.T, db DB) {
	t.Skip("(TODO) not implemented yet")
}

func TestDB_Stat(t *testing.T) {
	// Stat() (s Stat)

	t.Run("memory", func(t *testing.T) {
		testDBStat(t, NewMemoryDB())
	})

	t.Run("drive", func(t *testing.T) {
		db, cleanUp := testDriveDB(t)
		defer cleanUp()
		testDBStat(t, db)
	})

}

func testDBClose(t *testing.T, db DB) {
	if err := db.Close(); err != nil {
		t.Error("closing error:", err)
	}
	// Close can be called many times
	defer shouldNotPanic(t)
	db.Close()
}

func TestDB_Close(t *testing.T) {
	// Close() (err error)

	t.Run("memory", func(t *testing.T) {
		testDBClose(t, NewMemoryDB())
	})

	t.Run("drive", func(t *testing.T) {
		db, cleanUp := testDriveDB(t)
		defer cleanUp()
		testDBClose(t, db)
	})

}

//
// Tv
//

func testTvObjects(t *testing.T, db DB) {
	err := db.View(func(tx Tv) (_ error) {
		if tx.Objects() == nil {
			t.Error("Tv.Objects returns nil")
		}
		return
	})
	if err != nil {
		t.Error(err)
	}
	return
}

func TestTv_Objects(t *testing.T) {
	// Objects() ViewObjects

	t.Run("memory", func(t *testing.T) {
		testTvObjects(t, NewMemoryDB())
	})

	t.Run("drive", func(t *testing.T) {
		db, cleanUp := testDriveDB(t)
		defer cleanUp()
		testTvObjects(t, db)
	})

}

func testTvFeeds(t *testing.T, db DB) {
	err := db.View(func(tx Tv) (_ error) {
		if tx.Feeds() == nil {
			t.Error("Tv.Feeds returns nil")
		}
		return
	})
	if err != nil {
		t.Error(err)
	}
}

func TestTv_Feeds(t *testing.T) {
	// Feeds() ViewFeeds

	t.Run("memory", func(t *testing.T) {
		testTvFeeds(t, NewMemoryDB())
	})

	t.Run("drive", func(t *testing.T) {
		db, cleanUp := testDriveDB(t)
		defer cleanUp()
		testTvFeeds(t, db)
	})

}

//
// Tu
//

func testTuObjects(t *testing.T, db DB) {
	err := db.Update(func(tx Tu) (_ error) {
		if tx.Objects() == nil {
			t.Error("Tu.Objects returns nil")
		}
		return
	})
	if err != nil {
		t.Error(err)
	}
	return
}

func TestTu_Objects(t *testing.T) {
	// Objects() UpdateObjects

	t.Run("memory", func(t *testing.T) {
		testTuObjects(t, NewMemoryDB())
	})

	t.Run("drive", func(t *testing.T) {
		db, cleanUp := testDriveDB(t)
		defer cleanUp()
		testTuObjects(t, db)
	})

}

func testTuFeeds(t *testing.T, db DB) {
	err := db.Update(func(tx Tu) (_ error) {
		if tx.Feeds() == nil {
			t.Error("Tu.Feeds returns nil")
		}
		return
	})
	if err != nil {
		t.Error(err)
	}
}

func TestTu_Feeds(t *testing.T) {
	// Feeds() UpdateFeeds

	t.Run("memory", func(t *testing.T) {
		testTuFeeds(t, NewMemoryDB())
	})

	t.Run("drive", func(t *testing.T) {
		db, cleanUp := testDriveDB(t)
		defer cleanUp()
		testTuFeeds(t, db)
	})

}
