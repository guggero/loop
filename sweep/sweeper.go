package sweep

import (
	"context"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/lightninglabs/loop/lndclient"
	"github.com/lightninglabs/loop/swap"
	"github.com/lightningnetwork/lnd/input"
	"github.com/lightningnetwork/lnd/keychain"
	"github.com/lightningnetwork/lnd/lnwallet/chainfee"
	"github.com/lightningnetwork/lnd/sweep"
)

// Config is the set of configuration options that must be given to the sweeper
// to initialize it.
type Config struct {
	TxConfTarget        uint32
	BatchWindowDuration time.Duration
	SweeperStore        sweep.SweeperStore
}

// Sweeper creates htlc sweep txes.
type Sweeper struct {
	cfg *Config
	lnd *lndclient.LndServices

	*sweep.UtxoSweeper
}

func New(cfg *Config, lnd *lndclient.LndServices) *Sweeper {
	return &Sweeper{
		cfg: cfg,
		lnd: lnd,
		UtxoSweeper: sweep.New(&sweep.UtxoSweeperConfig{
			FeeEstimator: &feeEstimator{lnd},
			GenSweepScript: func() ([]byte, error) {
				return newSweepScript(lnd)
			},
			Signer: &signer{lnd},
			PublishTransaction: func(tx *wire.MsgTx) error {
				// TODO:context
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()

				return lnd.WalletKit.PublishTransaction(ctx, tx)
			},
			NewBatchTimer: func() <-chan time.Time {
				return time.NewTimer(cfg.BatchWindowDuration).C
			},
			Notifier:             &notifier{lnd},
			Store:                cfg.SweeperStore,
			MaxInputsPerTx:       sweep.DefaultMaxInputsPerTx,
			MaxSweepAttempts:     sweep.DefaultMaxSweepAttempts,
			NextAttemptDeltaFunc: sweep.DefaultNextAttemptDeltaFunc,
			MaxFeeRate:           sweep.DefaultMaxFeeRate,
			FeeRateBucketSize:    sweep.DefaultFeeRateBucketSize,
		}),
	}
}

func (s *Sweeper) CreateSweepTx(globalCtx context.Context, height uint32,
	htlc *swap.Htlc, htlcOutpoint wire.OutPoint, keyBytes [33]byte,
	amount btcutil.Amount, witnessType input.WitnessType) (
	*wire.MsgTx, error) {

	feePref := sweep.FeePreference{ConfTarget: s.cfg.TxConfTarget}
	return s.createSweepTx(
		globalCtx, height, htlc, htlcOutpoint, keyBytes, amount,
		witnessType, feePref,
	)
}

func (s *Sweeper) CreateSweepTxCustomFee(globalCtx context.Context,
	height uint32, htlc *swap.Htlc, htlcOutpoint wire.OutPoint,
	keyBytes [33]byte, amount btcutil.Amount, witnessType input.WitnessType,
	feeRate chainfee.SatPerKWeight) (*wire.MsgTx, error) {

	feePref := sweep.FeePreference{FeeRate: feeRate}
	return s.createSweepTx(
		globalCtx, height, htlc, htlcOutpoint, keyBytes, amount,
		witnessType, feePref,
	)
}

func (s *Sweeper) createSweepTx(globalCtx context.Context,
	height uint32, htlc *swap.Htlc, htlcOutpoint wire.OutPoint,
	keyBytes [33]byte, amount btcutil.Amount, witnessType input.WitnessType,
	feePref sweep.FeePreference) (*wire.MsgTx, error) {

	log.Infof("Publishing timeout tx")

	key, err := btcec.ParsePubKey(keyBytes[:], btcec.S256())
	if err != nil {
		return nil, err
	}

	var weightEstimate input.TxWeightEstimator
	weightEstimate.AddP2WKHOutput()
	htlc.AddTimeoutToEstimator(&weightEstimate)

	signDesc := &input.SignDescriptor{
		WitnessScript: htlc.Script,
		Output: &wire.TxOut{
			Value: int64(amount),
		},
		HashType:   txscript.SigHashAll,
		InputIndex: 0,
		KeyDesc: keychain.KeyDescriptor{
			PubKey: key,
		},
	}
	
	inp := input.MakeBaseInput(
		&htlcOutpoint, witnessType, signDesc, height,
	)

	// With our input constructed, we'll now offer it to the
	// sweeper.
	log.Infof("sweeping input")
	resultChan, err := s.SweepInput(&inp, feePref)
	if err != nil {
		return nil, err
	}

	log.Infof("waiting for result")
	// Sweeper is going to join this input with other inputs if
	// possible and publish the sweep tx. When the sweep tx
	// confirms, it signals us through the result channel with the
	// outcome. Wait for this to happen.
	select {
	case sweepResult := <-resultChan:
		if sweepResult.Err != nil {
			return nil, sweepResult.Err
		}

		log.Infof("Sweep success")
		return sweepResult.Tx, nil
	case <-globalCtx.Done():
		return nil, fmt.Errorf("quitting")
	}
}

// GetSweepFee calculates the required tx fee to spend to P2WKH. It takes a
// function that is expected to add the weight of the input to the weight
// estimator.
func (s *Sweeper) GetSweepFee(ctx context.Context,
	addInputEstimate func(*input.TxWeightEstimator),
	destAddr btcutil.Address, sweepConfTarget int32) (
	btcutil.Amount, int64, error) {

	// Get fee estimate from lnd.
	feeRate, err := s.lnd.WalletKit.EstimateFee(ctx, sweepConfTarget)
	if err != nil {
		return 0, 0, fmt.Errorf("estimate fee: %v", err)
	}

	// Calculate weight for this tx.
	var weightEstimate input.TxWeightEstimator
	switch destAddr.(type) {
	case *btcutil.AddressWitnessScriptHash:
		weightEstimate.AddP2WSHOutput()
	case *btcutil.AddressWitnessPubKeyHash:
		weightEstimate.AddP2WKHOutput()
	case *btcutil.AddressScriptHash:
		weightEstimate.AddP2SHOutput()
	case *btcutil.AddressPubKeyHash:
		weightEstimate.AddP2PKHOutput()
	default:
		return 0, 0, fmt.Errorf("unknown address type %T", destAddr)
	}

	addInputEstimate(&weightEstimate)
	weight := int64(weightEstimate.Weight())

	return feeRate.FeeForWeight(weight), weight, nil
}
