package main

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/MixinNetwork/mixin/common"
	"github.com/MixinNetwork/mixin/crypto"
	"github.com/MixinNetwork/mixin/storage"
	"github.com/dgraph-io/badger"
)

func decodeUtxoKey(key []byte) (crypto.Hash, int) {
	i, _ := binary.Varint(key[36:])
	return crypto.NewHash(key[4:36]), int(i)
}

func decodeUtxo(key, val []byte, pk crypto.Key) *common.Output {
	
	var out common.UTXOWithLock
	common.DecompressMsgpackUnmarshal(val, &out)
	
	_, i := decodeUtxoKey(key)
	return out.ViewGhostKey(&pk, i)

	// return &out.Output
}

func iterateUtxo() {
	s, err := storage.NewBadgerStore("/tmp/mixin-7001")

	db := *s.GetSnapshotsDB()

	err = db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte("UTXO")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {

				h, i := decodeUtxoKey(k)
				fmt.Printf("hash => %s, index => %d\n", h.String(), i)
				buf, _ := hex.DecodeString("9c7c8efae35fe0d464614d1ad9ae66923e6800814de5c4829015b927f239a807")
				var buf32 [32]byte
				copy(buf32[:], buf)

				out := decodeUtxo(k, v, crypto.Key(buf32))

				js, _ := json.MarshalIndent(out, "  ", " ")
				fmt.Printf("%s\n", string(js))


				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	fmt.Println(err)
}

func main() {

	iterateUtxo()
}