package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/boltdb/bolt"
	"github.com/compasses/MockXServer/utils"
)

type ReplayDB struct {
	db        *bolt.DB
	reqRspKey string
}

func NewReplayDB(path string) (*ReplayDB, error) {
	dbopen, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &ReplayDB{db: dbopen, reqRspKey: "|"}, nil
}

func (replay *ReplayDB) GetDBFilePath() string {
	return replay.db.Path()
}

func (replay *ReplayDB) StoreRequestFromJson(path, method string, reqBody, respBody interface{}, statusCode int) {
	req := utils.JsonInterfaceToByte(reqBody)
	rsp := utils.JsonInterfaceToByte(respBody)

	err := replay.StoreRequest(path, method, string(req), string(rsp), statusCode)
	if err != nil {
		log.Println("Store for json failed ", err)
	}
}

func (replay *ReplayDB) StoreRequest(path, method, reqBody, respBody string, statusCode int) (err error) {
	finalReq, finalResp, err := utils.JsonNormalize(reqBody, respBody, statusCode)
	if err != nil {
		log.Println("JSON Normalize error ", err)
	}

	replay.db.Update(func(tx *bolt.Tx) error {
		pathBucket, err := tx.CreateBucketIfNotExists([]byte(path))
		if err != nil {
			err = fmt.Errorf("create bucket: %s", err)
			return err
		}

		err = pathBucket.Put([]byte(method+replay.reqRspKey), finalResp)
		if err != nil {
			err = fmt.Errorf("store bucket error: %s", err)
			return err
		}

		methodBucket, err := pathBucket.CreateBucketIfNotExists([]byte(method))
		err = methodBucket.Put(finalReq, finalResp)
		if err != nil {
			err = fmt.Errorf("store bucket error: %s", err)
			return err
		}
		return nil
	})
	return
}

func (replay *ReplayDB) GetResponse(path, method, reqBody string) (resp []byte, err error) {
	replay.db.View(func(tx *bolt.Tx) error {

		pathBucket := tx.Bucket([]byte(path))
		if pathBucket == nil {
			err = fmt.Errorf("No response found for path:%s", path)
			return err
		}
		methodBucket := pathBucket.Bucket([]byte(method))
		if methodBucket == nil {
			err = fmt.Errorf("No response for path and method:%s", path+method)
			return err
		}
		var req []byte
		req, err = utils.JsonNormalizeSingle(reqBody)
		if err != nil {
			return err
		}
		resp = methodBucket.Get(req)
		if resp == nil {
			//use fuzzy match to get result
			resp = pathBucket.Get([]byte(method + replay.reqRspKey))
			if resp == nil {
				err = fmt.Errorf("No response for path and method:%s", path+method)
			}
		}
		return nil
	})
	return
}

func (replay *ReplayDB) Close() {
	err := replay.db.Close()
	if err != nil {
		fmt.Println("Close error:", err)
	}
}

func (replay *ReplayDB) Open(path string) {
	db, err := bolt.Open(path, 0600, nil)
	replay.db = db
	if err != nil {
		log.Println("Open DB error:", err)
	}
}

func (replay *ReplayDB) GetDB() *bolt.DB {
	return replay.db
}

func (replay *ReplayDB) GetJSONMap() (outmap map[string]map[string][]interface{}) {
	outmap = make(map[string]map[string][]interface{})

	replay.db.View(func(tx *bolt.Tx) error {
		cur := tx.Cursor()
		buc := cur.Bucket()
		c := buc.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			//fmt.Printf("Out side key=%s, value=%s\n", k, v)
			if v == nil {
				bb := tx.Bucket(k)
				if bb == nil {
					continue
				}
				cc := bb.Cursor()
				for kk, vv := cc.First(); kk != nil; kk, vv = cc.Next() {
					//fmt.Printf("middel layer key=%s, value=%s\n", kk, vv)
					if vv == nil {
						bbb := bb.Bucket(kk)
						if bbb == nil {
							continue
						}
						ccc := bbb.Cursor()
						for kkk, vvv := ccc.First(); kkk != nil; kkk, vvv = ccc.Next() {
							//fmt.Printf("last layer key=%s, value=%s\n", kkk, vvv)
							if outmap[string(k)] == nil {
								outmap[string(k)] = make(map[string][]interface{})
							}
							var pamv interface{}
							var resp interface{}

							err := json.Unmarshal(kkk, &pamv) //string(kkk)
							if err != nil {
								log.Println("Unmarshal failed ", err)
							}
							err = json.Unmarshal(vvv, &resp)
							if err != nil {
								log.Println("Unmarshal failed ", err)
							}
							outmap[string(k)][string(kk)] = append(outmap[string(k)][string(kk)], map[string]interface{}{
								"request":  pamv,
								"response": resp,
							})
						}
					}
				}
			}
		}
		return nil
	})
	return
}

func (replay *ReplayDB) SerilizeToFile() (filename string) {
	filename = "./single.json"
	outmap := replay.GetJSONMap()

	result := map[string]interface{}{
		"paths": outmap,
	}

	jsonStr, err := json.MarshalIndent(result, "", "    ")

	err = ioutil.WriteFile(filename, jsonStr, 0666)
	if err != nil {
		log.Println(err)
	}
	return
}
