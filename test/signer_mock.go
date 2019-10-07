package test

import (
	"context"

	"github.com/btcsuite/btcd/wire"
	"github.com/lightninglabs/loop/lndclient"
	"github.com/lightningnetwork/lnd/input"
)

type mockSigner struct {
}

var _ lndclient.SignerClient = (*mockSigner)(nil)

func (s *mockSigner) SignOutputRaw(ctx context.Context, tx *wire.MsgTx,
	signDescriptors []*input.SignDescriptor) ([][]byte, error) {

	rawSigs := [][]byte{{1, 2, 3}}

	return rawSigs, nil
}

func (s *mockSigner) ComputeInputScript(ctx context.Context, tx *wire.MsgTx,
	signDescriptors []*input.SignDescriptor) ([]*input.Script, error) {

	rawScripts := []*input.Script{{
		Witness:   [][]byte{{1, 2, 3}},
		SigScript: []byte{1, 2, 3},
	}}
	
	return rawScripts, nil
}
