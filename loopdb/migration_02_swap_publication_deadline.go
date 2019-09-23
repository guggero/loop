package loopdb

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/coreos/bbolt"
)

// deserializeLoopOutContractV01 is the version of loop out contract
// deserialization used _before_ v02 of the database, which is migrated to with
// the migration in this file.
func deserializeLoopOutContractV01(value []byte, chainParams *chaincfg.Params) (
	*LoopOutContract, error) {

	r := bytes.NewReader(value)

	contract := LoopOutContract{}
	var err error
	var unixNano int64
	if err := binary.Read(r, byteOrder, &unixNano); err != nil {
		return nil, err
	}
	contract.InitiationTime = time.Unix(0, unixNano)

	if err := binary.Read(r, byteOrder, &contract.Preimage); err != nil {
		return nil, err
	}

	err = binary.Read(r, byteOrder, &contract.AmountRequested)
	if err != nil {
		return nil, err
	}

	contract.PrepayInvoice, err = wire.ReadVarString(r, 0)
	if err != nil {
		return nil, err
	}

	n, err := r.Read(contract.SenderKey[:])
	if err != nil {
		return nil, err
	}
	if n != keyLength {
		return nil, fmt.Errorf("sender key has invalid length")
	}

	n, err = r.Read(contract.ReceiverKey[:])
	if err != nil {
		return nil, err
	}
	if n != keyLength {
		return nil, fmt.Errorf("receiver key has invalid length")
	}

	if err := binary.Read(r, byteOrder, &contract.CltvExpiry); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &contract.MaxMinerFee); err != nil {
		return nil, err
	}

	if err := binary.Read(r, byteOrder, &contract.MaxSwapFee); err != nil {
		return nil, err
	}

	if err := binary.Read(r, byteOrder, &contract.MaxPrepayRoutingFee); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &contract.InitiationHeight); err != nil {
		return nil, err
	}

	addr, err := wire.ReadVarString(r, 0)
	if err != nil {
		return nil, err
	}
	contract.DestAddr, err = btcutil.DecodeAddress(addr, chainParams)
	if err != nil {
		return nil, err
	}

	contract.SwapInvoice, err = wire.ReadVarString(r, 0)
	if err != nil {
		return nil, err
	}

	if err := binary.Read(r, byteOrder, &contract.SweepConfTarget); err != nil {
		return nil, err
	}

	if err := binary.Read(r, byteOrder, &contract.MaxSwapRoutingFee); err != nil {
		return nil, err
	}

	var unchargeChannel uint64
	if err := binary.Read(r, byteOrder, &unchargeChannel); err != nil {
		return nil, err
	}
	if unchargeChannel != 0 {
		contract.UnchargeChannel = &unchargeChannel
	}

	return &contract, nil
}

// migrateSwapPublicationDeadline migrates the database to v02, by adding the
// SwapPublicationDeadline field to loop out contracts.
func migrateSwapPublicationDeadline(tx *bbolt.Tx, chainParams *chaincfg.Params) error {
	rootBucket := tx.Bucket(loopOutBucketKey)
	if rootBucket == nil {
		return errors.New("bucket does not exist")
	}

	return rootBucket.ForEach(func(swapHash, v []byte) error {
		// Only go into things that we know are sub-bucket
		// keys.
		if v != nil {
			return nil
		}

		// From the root bucket, we'll grab the next swap
		// bucket for this swap from its swaphash.
		swapBucket := rootBucket.Bucket(swapHash)
		if swapBucket == nil {
			return fmt.Errorf("swap bucket %x not found",
				swapHash)
		}

		// With the main swap bucket obtained, we'll grab the
		// raw swap contract bytes and decode it.
		contractBytes := swapBucket.Get(contractKey)
		if contractBytes == nil {
			return errors.New("contract not found")
		}

		// Read the old format....
		contract, err := deserializeLoopOutContractV01(
			contractBytes, chainParams,
		)
		if err != nil {
			return err
		}

		// Since we don't have this field available for old swaps, set
		// the deadline to the initiation time, for immediate
		// publication.
		contract.SwapPublicationDeadline = time.Unix(
			contract.InitiationTime.Unix(), 0,
		)

		// ...write the new format.
		bs, err := serializeLoopOutContract(contract)
		if err != nil {
			return err
		}

		return swapBucket.Put(contractKey, bs)
	})
}
