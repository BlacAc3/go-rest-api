package database

import (
	"fmt"
	"log"
	"time"
    "os"
    "path/filepath"

	"github.com/boltdb/bolt"
)

var( 
    ProjectName = "go-rest-api"
    DBName string = "bolt_database.db"
    testDBName string = "test_bolt_database.db"
    Test_Mode bool = false

)


func FindProjectRoot(marker string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(cwd, marker)); err == nil {
			return cwd, nil
		}

		parent := filepath.Dir(cwd)
		if parent == cwd { // Reached root of filesystem
			break
		}
		cwd = parent
	}

	return "", fmt.Errorf("project root with marker %q not found", marker)
}



func OpenBoltDB() *bolt.DB{
    workingDir, err := FindProjectRoot("go.mod")
    if err != nil {
        return nil
    }
    if Test_Mode{
        boltDB, err := bolt.Open(fmt.Sprintf("%v/%v", workingDir, testDBName), 0600, &bolt.Options{Timeout: 1 * time.Second})
        if err != nil {
            log.Printf("Unable to open a Bolt DataBase:%v \n Found Error: %v", DBName, err)
            return nil
        }
        return boltDB
    }else{
        boltDB, err := bolt.Open(fmt.Sprintf("%v/%v", workingDir, DBName), 0600, &bolt.Options{Timeout: 1 * time.Second})
        if err != nil {
            log.Printf("Unable to open a Bolt DataBase:%v \n Found Error: %v", DBName, err)
            return nil
        }
        return boltDB
    }
}


func CreateBoltBucket(name string) error{
    db := OpenBoltDB()
    if db == nil {
        return fmt.Errorf("An Error occured while opening the database")
    }else{
        defer db.Close()
    }

    error := db.Update(func(tx *bolt.Tx) error {
        _, err := tx.CreateBucketIfNotExists([]byte(name))
        if err != nil {
            return err
        }
        return nil
    })
    return error
}


func UpdateBoltBucket(bucketName string, key string, value []byte) error{
    db := OpenBoltDB()
    if db == nil {
        return fmt.Errorf("An Error occured while opening the database")
    }else{
        defer db.Close()
    }

    error := db.Update(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(bucketName))
        if b == nil {
            return nil
        }
        return b.Put([]byte(key), value)
    })
    return error
}


func GetBoltBucket(name string, key string) ([]byte ,error){
    db := OpenBoltDB()
    if db == nil {
        return nil, fmt.Errorf("An Error occured while getting the bucket: %v", name)
    }else{
        defer db.Close()
    }
    var data string
    err := db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte(name))
        if b == nil {
            return fmt.Errorf("An error occured while getting the bucket: %v", name)
        }
        data = string(b.Get([]byte(key)))
        if data == ""{
            return fmt.Errorf("User not Found!")
        }
        return nil
    })
    if err != nil{
        return nil, err
    }
    return []byte(data), nil

}
