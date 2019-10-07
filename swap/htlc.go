package swap

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/lightningnetwork/lnd/input"
	"github.com/lightningnetwork/lnd/lntypes"
)

// HtlcOutputType defines the output type of the htlc that is published.
type HtlcOutputType uint8

type HtlcSweepType uint8

const (
	// HtlcP2WSH is a pay-to-witness-script-hash output (segwit only)
	HtlcP2WSH HtlcOutputType = iota

	// HtlcNP2WSH is a nested pay-to-witness-script-hash output that can be
	// paid to be legacy wallets.
	HtlcNP2WSH

	HtlcSuccess HtlcSweepType = 0

	HtlcTimeout HtlcSweepType = 1
)

// Htlc contains relevant htlc information from the receiver perspective.
type Htlc struct {
	Script      []byte
	PkScript    []byte
	Hash        lntypes.Hash
	OutputType  HtlcOutputType
	ChainParams *chaincfg.Params
	Address     btcutil.Address
	SigScript   []byte
}

var (
	quoteKey [33]byte

	quoteHash lntypes.Hash

	// QuoteHtlc is a template script just used for fee estimation. It uses
	// the maximum value for cltv expiry to get the maximum (worst case)
	// script size.
	QuoteHtlc, _ = NewHtlc(
		^int32(0), quoteKey, quoteKey, quoteHash, HtlcP2WSH,
		&chaincfg.MainNetParams,
	)
)

// NewHtlc returns a new instance.
func NewHtlc(cltvExpiry int32, senderKey, receiverKey [33]byte,
	hash lntypes.Hash, outputType HtlcOutputType,
	chainParams *chaincfg.Params) (*Htlc, error) {

	script, err := swapHTLCScript(
		cltvExpiry, senderKey, receiverKey, hash,
	)
	if err != nil {
		return nil, err
	}

	p2wshPkScript, err := input.WitnessScriptHash(script)
	if err != nil {
		return nil, err
	}

	var pkScript, sigScript []byte
	var address btcutil.Address

	switch outputType {
	case HtlcNP2WSH:
		// Generate p2sh script for p2wsh (nested).
		p2wshPkScriptHash := sha256.Sum256(p2wshPkScript)
		hash160 := input.Ripemd160H(p2wshPkScriptHash[:])

		builder := txscript.NewScriptBuilder()

		builder.AddOp(txscript.OP_HASH160)
		builder.AddData(hash160)
		builder.AddOp(txscript.OP_EQUAL)

		pkScript, err = builder.Script()
		if err != nil {
			return nil, err
		}

		// Generate a valid sigScript that will allow us to spend the
		// p2sh output. The sigScript will contain only a single push of
		// the p2wsh witness program corresponding to the matching
		// public key of this address.
		sigScript, err = txscript.NewScriptBuilder().
			AddData(p2wshPkScript).
			Script()
		if err != nil {
			return nil, err
		}

		address, err = btcutil.NewAddressScriptHash(
			p2wshPkScript, chainParams,
		)
		if err != nil {
			return nil, err
		}

	case HtlcP2WSH:
		pkScript = p2wshPkScript

		address, err = btcutil.NewAddressWitnessScriptHash(
			p2wshPkScript[2:],
			chainParams,
		)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unknown output type")
	}

	return &Htlc{
		Hash:        hash,
		Script:      script,
		PkScript:    pkScript,
		OutputType:  outputType,
		ChainParams: chainParams,
		Address:     address,
		SigScript:   sigScript,
	}, nil
}

// SwapHTLCScript returns the on-chain HTLC witness script.
//
// OP_SIZE 32 OP_EQUAL
// OP_IF
//    OP_HASH160 <ripemd160(swap_hash)> OP_EQUALVERIFY
//    <recvr key>
// OP_ELSE
//    OP_DROP
//    <cltv timeout> OP_CHECKLOCKTIMEVERIFY OP_DROP
//    <sender key>
// OP_ENDIF
// OP_CHECKSIG
func swapHTLCScript(cltvExpiry int32, senderHtlcKey,
	receiverHtlcKey [33]byte, swapHash lntypes.Hash) ([]byte, error) {

	builder := txscript.NewScriptBuilder()

	builder.AddOp(txscript.OP_SIZE)
	builder.AddInt64(32)
	builder.AddOp(txscript.OP_EQUAL)

	builder.AddOp(txscript.OP_IF)

	builder.AddOp(txscript.OP_HASH160)
	builder.AddData(input.Ripemd160H(swapHash[:]))
	builder.AddOp(txscript.OP_EQUALVERIFY)

	builder.AddData(receiverHtlcKey[:])

	builder.AddOp(txscript.OP_ELSE)

	builder.AddOp(txscript.OP_DROP)

	builder.AddInt64(int64(cltvExpiry))
	builder.AddOp(txscript.OP_CHECKLOCKTIMEVERIFY)
	builder.AddOp(txscript.OP_DROP)

	builder.AddData(senderHtlcKey[:])

	builder.AddOp(txscript.OP_ENDIF)

	builder.AddOp(txscript.OP_CHECKSIG)

	return builder.Script()
}

// GenSuccessWitness returns the success script to spend this htlc with the
// preimage.
func (h *Htlc) GenSuccessWitness(receiverSig []byte,
	preimage lntypes.Preimage) (wire.TxWitness, error) {

	if h.Hash != preimage.Hash() {
		return nil, errors.New("preimage doesn't match hash")
	}

	witnessStack := make(wire.TxWitness, 3)
	witnessStack[0] = append(receiverSig, byte(txscript.SigHashAll))
	witnessStack[1] = preimage[:]
	witnessStack[2] = h.Script

	return witnessStack, nil
}

// IsSuccessWitness checks whether the given stack is valid for redeeming the
// htlc.
func (h *Htlc) IsSuccessWitness(witness wire.TxWitness) bool {
	if len(witness) != 3 {
		return false
	}

	isTimeoutTx := bytes.Equal([]byte{0}, witness[1])

	return !isTimeoutTx
}

// GenTimeoutWitness returns the timeout script to spend this htlc after
// timeout.
func (h *Htlc) GenTimeoutWitness(senderSig []byte) (wire.TxWitness, error) {

	witnessStack := make(wire.TxWitness, 3)
	witnessStack[0] = append(senderSig, byte(txscript.SigHashAll))
	witnessStack[1] = []byte{0}
	witnessStack[2] = h.Script

	return witnessStack, nil
}

func (h *Htlc) MaxSuccessWitnessSize() int {
	// Calculate maximum success witness size
	//
	// - number_of_witness_elements: 1 byte
	// - receiver_sig_length: 1 byte
	// - receiver_sig: 73 bytes
	// - preimage_length: 1 byte
	// - preimage: 33 bytes
	// - witness_script_length: 1 byte
	// - witness_script: len(script) bytes
	return 1 + 1 + 73 + 1 + 33 + 1 + len(h.Script)
}

// AddSuccessToEstimator adds a successful spend to a weight estimator.
func (h *Htlc) AddSuccessToEstimator(estimator *input.TxWeightEstimator) {
	switch h.OutputType {
	case HtlcP2WSH:
		estimator.AddWitnessInput(h.MaxSuccessWitnessSize())

	case HtlcNP2WSH:
		estimator.AddNestedP2WSHInput(h.MaxSuccessWitnessSize())
	}
}

func (h *Htlc) MaxTimeoutWitnessSize() int {
	// Calculate maximum timeout witness size
	//
	// - number_of_witness_elements: 1 byte
	// - sender_sig_length: 1 byte
	// - sender_sig: 73 bytes
	// - zero_length: 1 byte
	// - zero: 1 byte
	// - witness_script_length: 1 byte
	// - witness_script: len(script) bytes
	return 1 + 1 + 73 + 1 + 1 + 1 + len(h.Script)
}

// AddTimeoutToEstimator adds a timeout spend to a weight estimator.
func (h *Htlc) AddTimeoutToEstimator(estimator *input.TxWeightEstimator) {
	switch h.OutputType {
	case HtlcP2WSH:
		estimator.AddWitnessInput(h.MaxTimeoutWitnessSize())

	case HtlcNP2WSH:
		estimator.AddNestedP2WSHInput(h.MaxTimeoutWitnessSize())
	}
}

type SweepWitness struct {
	htlc          *Htlc
	sweepType     HtlcSweepType
	preimage      lntypes.Preimage
	currentHeight uint32
}

var _ input.WitnessType = (*SweepWitness)(nil)

func NewSuccessSweepWitness(htlc *Htlc, currentHeight uint32,
	preimage lntypes.Preimage) *SweepWitness {

	return &SweepWitness{
		htlc:          htlc,
		currentHeight: currentHeight,
		sweepType:     HtlcSuccess,
		preimage:      preimage,
	}
}

func NewTimeoutSweepWitness(htlc *Htlc, currentHeight uint32) *SweepWitness {
	return &SweepWitness{
		htlc:          htlc,
		currentHeight: currentHeight,
		sweepType:     HtlcTimeout,
	}
}

func (sw *SweepWitness) String() string {
	switch sw.sweepType {
	case HtlcSuccess:
		return "LoopHtlcSuccess"
	case HtlcTimeout:
		return "LoopHtlcTimeout"
	default:
		return fmt.Sprintf(
			"Unknown WitnessType: %v", uint8(sw.sweepType),
		)
	}
}

func (sw *SweepWitness) WitnessGenerator(signer input.Signer,
	descriptor *input.SignDescriptor) input.WitnessGenerator {

	return func(tx *wire.MsgTx, hc *txscript.TxSigHashes,
		inputIndex int) (*input.Script, error) {

		desc := descriptor
		desc.SigHashes = hc
		desc.InputIndex = inputIndex

		tx.LockTime = sw.currentHeight

		if sw.htlc.OutputType == HtlcNP2WSH {
			tx.TxIn[inputIndex].SignatureScript = sw.htlc.SigScript
		}

		rawSig, err := signer.SignOutputRaw(
			tx, descriptor,
		)
		if err != nil {
			return nil, fmt.Errorf("could not sign: %v", err)
		}

		switch sw.sweepType {
		case HtlcSuccess:
			// Add witness stack to the tx input.
			witness, err := sw.htlc.GenSuccessWitness(
				rawSig, sw.preimage,
			)
			if err != nil {
				return nil, err
			}

			return &input.Script{
				Witness: witness,
			}, nil
		case HtlcTimeout:
			// Add witness stack to the tx input.
			witness, err := sw.htlc.GenTimeoutWitness(rawSig)
			if err != nil {
				return nil, err
			}

			return &input.Script{
				Witness: witness,
			}, nil
		default:
			return nil, fmt.Errorf(
				"unknown sweep type: %v", sw.sweepType,
			)
		}
	}
}

func (sw *SweepWitness) SizeUpperBound() (int, bool, error) {
	isNestedP2SH := sw.htlc.OutputType == HtlcNP2WSH

	switch sw.sweepType {
	case HtlcSuccess:
		return sw.htlc.MaxSuccessWitnessSize(), isNestedP2SH, nil
	case HtlcTimeout:
		return sw.htlc.MaxTimeoutWitnessSize(), isNestedP2SH, nil
	default:
		return 0, false, fmt.Errorf(
			"unknown sweep type: %v", sw.sweepType,
		)
	}
}
func (sw *SweepWitness) AddWeightEstimation(
	estimator *input.TxWeightEstimator) error {

	switch sw.sweepType {
	case HtlcSuccess:
		sw.htlc.AddSuccessToEstimator(estimator)
	case HtlcTimeout:
		sw.htlc.AddSuccessToEstimator(estimator)
	default:
		return fmt.Errorf("unknown sweep type: %v", sw.sweepType)
	}

	return nil
}
