package test

import (
	"fmt"
	"sync"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/lightninglabs/loop/lndclient"
	"github.com/lightningnetwork/lnd/chainntnfs"
	"golang.org/x/net/context"
)

type mockChainNotifier struct {
	lnd         *LndMockServices
	wg          sync.WaitGroup
	blockHash   *chainhash.Hash
	blockHeight int32
}

var _ lndclient.ChainNotifierClient = (*mockChainNotifier)(nil)

// SpendRegistration contains registration details.
type SpendRegistration struct {
	Outpoint   *wire.OutPoint
	PkScript   []byte
	HeightHint int32
}

// ConfRegistration contains registration details.
type ConfRegistration struct {
	TxID       *chainhash.Hash
	PkScript   []byte
	HeightHint int32
	NumConfs   int32
}

func (c *mockChainNotifier) RegisterSpendNtfn(ctx context.Context,
	outpoint *wire.OutPoint, pkScript []byte, heightHint int32) (
	chan *chainntnfs.SpendDetail, chan error, error) {

	fmt.Println("register spend chan")
	c.lnd.RegisterSpendChannel <- &SpendRegistration{
		HeightHint: heightHint,
		Outpoint:   outpoint,
		PkScript:   pkScript,
	}
	fmt.Println("sent on register spend chan")

	spendChan := make(chan *chainntnfs.SpendDetail, 1)
	errChan := make(chan error, 1)

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		select {
		case m := <-c.lnd.SpendChannel:
			select {
			case spendChan <- m:
			case <-ctx.Done():
			}
		case <-ctx.Done():
		}
	}()

	return spendChan, errChan, nil
}

func (c *mockChainNotifier) WaitForFinished() {
	c.wg.Wait()
}

func (c *mockChainNotifier) RegisterBlockEpochNtfn(ctx context.Context,
	bestBlock *chainntnfs.BlockEpoch) (chan *chainntnfs.BlockEpoch,
	chan error, error) {

	blockErrorChan := make(chan error, 1)
	blockEpochChan := make(chan *chainntnfs.BlockEpoch)

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		initialEpoch := &chainntnfs.BlockEpoch{Height: c.lnd.Height}

		// Send initial block height
		select {
		case blockEpochChan <- initialEpoch:
		case <-ctx.Done():
			return
		}

		for {
			select {
			case m := <-c.lnd.epochChannel:

				epoch := &chainntnfs.BlockEpoch{
					Height: m,
				}
				select {
				case blockEpochChan <- epoch:
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return blockEpochChan, blockErrorChan, nil
}

func (c *mockChainNotifier) RegisterConfirmationsNtfn(ctx context.Context,
	txid *chainhash.Hash, pkScript []byte, numConfs, heightHint int32) (
	chan *chainntnfs.TxConfirmation, chan error, error) {

	confChan := make(chan *chainntnfs.TxConfirmation, 1)
	errChan := make(chan error, 1)

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		select {
		case m := <-c.lnd.ConfChannel:
			select {
			case confChan <- m:
			case <-ctx.Done():
			}
		case <-ctx.Done():
		}
	}()

	select {
	case c.lnd.RegisterConfChannel <- &ConfRegistration{
		PkScript:   pkScript,
		TxID:       txid,
		HeightHint: heightHint,
		NumConfs:   numConfs,
	}:
	case <-time.After(Timeout):
		return nil, nil, ErrTimeout
	}

	return confChan, errChan, nil
}
