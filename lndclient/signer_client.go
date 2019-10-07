package lndclient

import (
	"context"

	"github.com/btcsuite/btcd/wire"
	"github.com/lightninglabs/loop/swap"
	"github.com/lightningnetwork/lnd/input"
	"github.com/lightningnetwork/lnd/lnrpc/signrpc"
	"google.golang.org/grpc"
)

// SignerClient exposes sign functionality.
type SignerClient interface {
	SignOutputRaw(ctx context.Context, tx *wire.MsgTx,
		signDescriptors []*input.SignDescriptor) ([][]byte, error)

	ComputeInputScript(ctx context.Context, tx *wire.MsgTx,
		signDescriptors []*input.SignDescriptor) (
		[]*input.Script, error)
}

type signerClient struct {
	client    signrpc.SignerClient
	signerMac serializedMacaroon
}

var _ SignerClient = (*signerClient)(nil)

func newSignerClient(conn *grpc.ClientConn,
	signerMac serializedMacaroon) *signerClient {

	return &signerClient{
		client:    signrpc.NewSignerClient(conn),
		signerMac: signerMac,
	}
}

func convertDescriptors(
	signDescriptors []*input.SignDescriptor) []*signrpc.SignDescriptor {

	rpcSignDescriptors := make(
		[]*signrpc.SignDescriptor, len(signDescriptors),
	)
	for i, signDesc := range signDescriptors {
		var keyBytes []byte
		var keyLocator *signrpc.KeyLocator
		if signDesc.KeyDesc.PubKey != nil {
			keyBytes = signDesc.KeyDesc.PubKey.SerializeCompressed()
		} else {
			keyLocator = &signrpc.KeyLocator{
				KeyFamily: int32(
					signDesc.KeyDesc.KeyLocator.Family,
				),
				KeyIndex: int32(
					signDesc.KeyDesc.KeyLocator.Index,
				),
			}
		}

		var doubleTweak []byte
		if signDesc.DoubleTweak != nil {
			doubleTweak = signDesc.DoubleTweak.Serialize()
		}

		rpcSignDescriptors[i] = &signrpc.SignDescriptor{
			WitnessScript: signDesc.WitnessScript,
			Output: &signrpc.TxOut{
				PkScript: signDesc.Output.PkScript,
				Value:    signDesc.Output.Value,
			},
			Sighash:    uint32(signDesc.HashType),
			InputIndex: int32(signDesc.InputIndex),
			KeyDesc: &signrpc.KeyDescriptor{
				RawKeyBytes: keyBytes,
				KeyLoc:      keyLocator,
			},
			SingleTweak: signDesc.SingleTweak,
			DoubleTweak: doubleTweak,
		}
	}
	return rpcSignDescriptors
}

func (s *signerClient) SignOutputRaw(ctx context.Context, tx *wire.MsgTx,
	signDescriptors []*input.SignDescriptor) ([][]byte, error) {

	txRaw, err := swap.EncodeTx(tx)
	if err != nil {
		return nil, err
	}

	rpcSignDescriptors := convertDescriptors(signDescriptors)

	rpcCtx, cancel := context.WithTimeout(ctx, rpcTimeout)
	defer cancel()

	rpcCtx = s.signerMac.WithMacaroonAuth(rpcCtx)
	resp, err := s.client.SignOutputRaw(rpcCtx,
		&signrpc.SignReq{
			RawTxBytes: txRaw,
			SignDescs:  rpcSignDescriptors,
		},
	)
	if err != nil {
		return nil, err
	}

	return resp.RawSigs, nil
}

func (s *signerClient) ComputeInputScript(ctx context.Context, tx *wire.MsgTx,
	signDescriptors []*input.SignDescriptor) ([]*input.Script, error) {

	txRaw, err := swap.EncodeTx(tx)
	if err != nil {
		return nil, err
	}

	rpcSignDescriptors := convertDescriptors(signDescriptors)

	rpcCtx, cancel := context.WithTimeout(ctx, rpcTimeout)
	defer cancel()

	rpcCtx = s.signerMac.WithMacaroonAuth(rpcCtx)
	resp, err := s.client.ComputeInputScript(rpcCtx,
		&signrpc.SignReq{
			RawTxBytes: txRaw,
			SignDescs:  rpcSignDescriptors,
		},
	)
	if err != nil {
		return nil, err
	}

	scripts := make([]*input.Script, len(resp.InputScripts))
	for i, script := range resp.InputScripts {
		scripts[i] = &input.Script{
			Witness:   script.Witness,
			SigScript: script.SigScript,
		}
	}

	return scripts, nil
}
