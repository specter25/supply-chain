package blockchain

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/dgraph-io/badger"
)

const (
	dbPath      = "./tmp/blocks"          // this %s allows us to have multiple databases for each of our nodes
	dbFile      = "./tmp/blocks/MANIFEST" //tocheck whether database exists or not
	genesisData = "First Transaction from genesis"
)

type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}

//helper function to check whether blockchain exists or not
func DBexisits(path string) bool {
	if _, err := os.Stat(path + "/MANIFEST"); os.IsNotExist(err) {
		return false
	}
	return true
}

//if one instance of the datatabse is already running then it
// will conatin the lock file and if the instance terminates without properly garbage removal
//then this lock file will remain in the folder

//retry function deletes the badger file and change the original optio

func retry(dir string, originalOpts badger.Options) (*badger.DB, error) {
	lockPath := filepath.Join(dir, "LOCK")
	if err := os.Remove(lockPath); err != nil {
		return nil, fmt.Errorf(`removing "LOCK": %s`, err)
	}
	retryOpts := originalOpts
	retryOpts.Truncate = true
	db, err := badger.Open(retryOpts)
	return db, err
}

func openDB(dir string, opts badger.Options) (*badger.DB, error) {
	if db, err := badger.Open(opts); err != nil {
		if strings.Contains(err.Error(), "LOCK") {
			if db, err := retry(dir, opts); err == nil {
				log.Println("database unlocked, value log truncated")
				return db, nil
			}
			log.Println("could not unlock database:", err)
		}
		return nil, err
	} else {
		return db, nil
	}
}

func InitBlockchain(address string, nodeId string) *Blockchain {
	path := fmt.Sprintf(dbPath, nodeId)
	var lastHash []byte

	if DBexisits(path) {
		fmt.Println("Blockchain does not exists")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(path)
	opts.Dir = path      //store keys and metadata
	opts.ValueDir = path //databse will store all the values

	db, err := openDB(path, opts)
	Handle(err)
	err = db.Update(func(txn *badger.Txn) error {
		//check if there is a blockchain already stored or not
		//if there is alreadya blockchainthen we will create a new blockchain instance in memory and we will get the
		//last hash of our blockchian in our disk database and we will push to this instance in memory
		//the reason why the last hash isimportant is that it helps derive a new block in our blockchain
		//if there is no existing blockchain we will create a genesis block we will push it in our databse then we will save the genesis block hash as the lastblock hash in our databse
		genesis := Genesis()
		fmt.Println("Genesis created")
		err = txn.Set(genesis.Hash, genesis.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), genesis.Hash)
		lastHash = genesis.Hash
		return err
	})

	Handle(err)
	blockchain := Blockchain{lastHash, db} // create a new blockchain in the memory
	return &blockchain                     // return a referenec to this blockchain

}
func ContinueBlockchain(nodeId string) *Blockchain {
	path := fmt.Sprintf(dbPath, nodeId)
	if DBexisits(path) == false {
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}
	var lastHash []byte
	opts := badger.DefaultOptions(path)
	opts.Dir = path      //store keys and metadata
	opts.ValueDir = path //databse will store all the values

	db, err := openDB(path, opts)
	Handle(err)
	err = db.Update(func(txn *badger.Txn) error {

		item, err := txn.Get([]byte("lh"))
		Handle(err)
		err1 := item.Value(func(val []byte) error {
			lastHash = append([]byte{}, val...)
			return nil
		})
		Handle(err1)
		return err
	})
	Handle(err)
	chain := Blockchain{lastHash, db}
	return &chain

}
