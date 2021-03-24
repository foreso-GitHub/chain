package node

import (
	"fmt"
	"sync"

	"github.com/tokentransfer/go-MerklePatriciaTree/mpt"

	"github.com/tokentransfer/chain/block"
	"github.com/tokentransfer/chain/crypto"
	"github.com/tokentransfer/chain/store"

	libblock "github.com/tokentransfer/interfaces/block"
	libcore "github.com/tokentransfer/interfaces/core"
	libcrypto "github.com/tokentransfer/interfaces/crypto"
	libnode "github.com/tokentransfer/interfaces/node"
	libstore "github.com/tokentransfer/interfaces/store"
)

type MerkleTree struct {
	mt     *mpt.Trie
	locker *sync.RWMutex
}

func NewMerkleTree(cs libcrypto.CryptoService, ss libstore.KvService) *MerkleTree {
	mt := mpt.New(cs, ss)
	return &MerkleTree{
		mt:     mt,
		locker: &sync.RWMutex{},
	}
}

func (t *MerkleTree) GetRoot() []byte {
	t.locker.RLock()
	defer t.locker.RUnlock()

	return t.mt.RootHash()
}

func (t *MerkleTree) Commit() error {
	t.locker.Lock()
	defer t.locker.Unlock()

	return t.mt.Commit()
}

func (t *MerkleTree) Cancel() error {
	t.locker.Lock()
	defer t.locker.Unlock()

	return t.mt.Abort()
}

func (t *MerkleTree) GetData(key []byte) ([]byte, error) {
	t.locker.RLock()
	defer t.locker.RUnlock()

	return t.mt.Get(key)
}

func (t *MerkleTree) PutData(key, value []byte) error {
	t.locker.Lock()
	defer t.locker.Unlock()

	return t.mt.Put(key, value)
}

type MerkleService struct {
	config libcore.Config

	im libnode.MerkleTree // index -> hash
	bm libnode.MerkleTree // block
	tm libnode.MerkleTree // transaction
	sm libnode.MerkleTree // state

	CryptoService *crypto.CryptoService
}

func (service *MerkleService) Init(c libcore.Config) error {
	service.config = c

	indexdb := &store.LevelService{Name: "index"}
	err := indexdb.Init(c)
	if err != nil {
		return err
	}
	service.im = NewMerkleTree(service.CryptoService, indexdb)

	blockdb := &store.LevelService{Name: "block"}
	err = blockdb.Init(c)
	if err != nil {
		return err
	}
	service.bm = NewMerkleTree(service.CryptoService, blockdb)

	txdb := &store.LevelService{Name: "transaction"}
	err = txdb.Init(c)
	if err != nil {
		return err
	}
	service.tm = NewMerkleTree(service.CryptoService, txdb)

	statedb := &store.LevelService{Name: "state"}
	err = statedb.Init(c)
	if err != nil {
		return err
	}
	service.sm = NewMerkleTree(service.CryptoService, statedb)
	return nil
}

func (service *MerkleService) Start() error {
	return nil
}

func (service *MerkleService) Close() error {
	return nil
}

func (service *MerkleService) PutState(s libblock.State) error {
	cs := service.CryptoService

	h, data, err := cs.Raw(s, libcrypto.RawBinary)
	if err != nil {
		return err
	}
	err = service.sm.PutData(h, data)
	if err != nil {
		return err
	}

	key := s.GetStateKey()
	indexKey := getIndexKey(key, s.GetIndex())
	stateKey := getNameKey("state", indexKey)
	err = service.im.PutData([]byte(stateKey), h)
	if err != nil {
		return err
	}

	newKey := getNameKey("state", key)
	err = service.im.PutData([]byte(newKey), h)
	if err != nil {
		return err
	}
	return nil
}

func (service *MerkleService) GetState(h libcore.Hash) (libblock.State, error) {
	data, err := service.sm.GetData(h)
	if err != nil {
		return nil, err
	}
	state, err := block.ReadState(data)
	if err != nil {
		return nil, err
	}
	return state, nil
}

func (service *MerkleService) GetStateByIndex(key string, index uint64) (libblock.State, error) {
	indexKey := getIndexKey(key, index)
	stateKey := getNameKey("state", indexKey)
	h, err := service.im.GetData([]byte(stateKey))
	if err != nil {
		return nil, err
	}
	return service.GetState(libcore.Hash(h))
}

func (service *MerkleService) GetStateByKey(key string) (libblock.State, error) {
	newKey := getNameKey("state", key)
	h, err := service.im.GetData([]byte(newKey))
	if err != nil {
		return nil, err
	}
	return service.GetState(libcore.Hash(h))
}

func (service *MerkleService) GetStateRoot() libcore.Hash {
	return service.sm.GetRoot()
}

func (service *MerkleService) PutTransaction(txWithData libblock.TransactionWithData) error {
	cs := service.CryptoService

	h, data, err := cs.Raw(txWithData, libcrypto.RawBinary)
	if err != nil {
		return err
	}
	err = service.tm.PutData(h, data)
	if err != nil {
		return err
	}

	txHash, _, err := cs.Raw(txWithData.GetTransaction(), libcrypto.RawBinary)
	txKey := getHashKey("transaction", txHash)
	err = service.im.PutData([]byte(txKey), h)
	if err != nil {
		return err
	}

	account := txWithData.GetTransaction().GetAccount()
	address, err := account.GetAddress()
	if err != nil {
		return err
	}
	indexKey := getIndexKey(address, txWithData.GetTransaction().GetIndex())
	accountKey := getNameKey("transaction", indexKey)
	err = service.im.PutData([]byte(accountKey), h)
	if err != nil {
		return err
	}

	return nil
}

func (service *MerkleService) GetTransaction(h libcore.Hash) (libblock.TransactionWithData, error) {
	data, err := service.tm.GetData(h)
	if err != nil {
		return nil, err
	}
	//txWithData := &block.TransactionWithData{}
	//err = txWithData.UnmarshalBinary(data)
	txWithData, err := block.ReadTxWithData(data)
	if err != nil {
		return nil, err
	}
	return txWithData, nil
}

func (service *MerkleService) GetTransactionByHash(txHash libcore.Hash) (libblock.TransactionWithData, error) {
	txKey := getHashKey("transaction", txHash)
	h, err := service.im.GetData([]byte(txKey))
	if err != nil {
		return nil, err
	}
	return service.GetTransaction(libcore.Hash(h))
}

func (service *MerkleService) GetTransactionByIndex(account libcore.Address, index uint64) (libblock.TransactionWithData, error) {
	address, err := account.GetAddress()
	if err != nil {
		return nil, err
	}
	indexKey := getIndexKey(address, index)
	accountKey := getNameKey("transaction", indexKey)
	h, err := service.im.GetData([]byte(accountKey))
	if err != nil {
		return nil, err
	}
	return service.GetTransaction(libcore.Hash(h))
}

func (service *MerkleService) GetTransactionRoot() libcore.Hash {
	return service.tm.GetRoot()
}

func (service *MerkleService) PutBlock(b libblock.Block) error {
	cs := service.CryptoService

	h, data, err := cs.Raw(b, libcrypto.RawBinary)
	if err != nil {
		return err
	}
	err = service.bm.PutData(h, data)
	if err != nil {
		return err
	}
	name := getBlockKey(b.GetIndex())
	err = service.im.PutData([]byte(name), h[:])
	if err != nil {
		return err
	}

	transactions := b.GetTransactions()
	l := len(transactions)
	for i := 0; i < l; i++ {
		tx := transactions[i]
		err := service.PutTransaction(tx)
		if err != nil {
			return err
		}
	}

	states := b.GetStates()
	l = len(states)
	for i := 0; i < l; i++ {
		s := states[i]
		err := service.PutState(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (service *MerkleService) GetBlockByIndex(index uint64) (libblock.Block, error) {
	name := getBlockKey(index)
	data, err := service.im.GetData([]byte(name))
	if err != nil {
		return nil, err
	}
	h := libcore.Hash(data)
	return service.GetBlockByHash(h)
}

func (service *MerkleService) GetBlockByHash(hash libcore.Hash) (libblock.Block, error) {
	data, err := service.bm.GetData(hash)
	if err != nil {
		return nil, err
	}
	b := &block.Block{}
	err = b.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (service *MerkleService) Commit() error {
	err := service.im.Commit()
	if err != nil {
		return err
	}
	err = service.bm.Commit()
	if err != nil {
		return err
	}
	err = service.tm.Commit()
	if err != nil {
		return err
	}
	err = service.sm.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (service *MerkleService) Cancel() error {
	err := service.im.Cancel()
	if err != nil {
		return err
	}
	err = service.bm.Cancel()
	if err != nil {
		return err
	}
	err = service.tm.Cancel()
	if err != nil {
		return err
	}
	err = service.sm.Cancel()
	if err != nil {
		return err
	}
	return nil
}

func getBlockKey(index uint64) string {
	return fmt.Sprintf("block@%d", index)
}

func getHashKey(name string, h libcore.Hash) string {
	return fmt.Sprintf("%s@%s", name, h.String())
}

func getNameKey(name string, s string) string {
	return fmt.Sprintf("%s@%s", name, s)
}

func getIndexKey(key string, index uint64) string {
	return fmt.Sprintf("%s:%d", key, index)
}
