package block

import (
	"errors"

	"github.com/tokentransfer/chain/account"
	"github.com/tokentransfer/chain/core"
	"github.com/tokentransfer/chain/core/pb"

	libblock "github.com/tokentransfer/interfaces/block"
	libcore "github.com/tokentransfer/interfaces/core"
)

type State struct {
	Hash       libcore.Hash
	BlockIndex uint64

	StateType libblock.StateType
}

func (s *State) GetHash() libcore.Hash {
	return s.Hash
}

func (s *State) SetHash(h libcore.Hash) {
	s.Hash = h
}

func (s *State) GetStateType() libblock.StateType {
	return s.StateType
}

func (s *State) GetBlockIndex() uint64 {
	return s.BlockIndex
}

func (s *State) SetBlockIndex(index uint64) {
	s.BlockIndex = index
}

type AccountState struct {
	State

	Account  libcore.Address
	Sequence uint64
	Amount   int64
}

func (s *AccountState) GetIndex() uint64 {
	return s.Sequence
}

func (s *AccountState) GetStateKey() string {
	a, err := s.Account.GetAddress()
	if err != nil {
		return ""
	}
	return a
}

func (s *AccountState) UnmarshalBinary(data []byte) error {
	meta, msg, err := core.Unmarshal(data)
	if err != nil {
		return err
	}
	if meta != core.CORE_ACCOUNT_STATE {
		return errors.New("error State data")
	}

	state := msg.(*pb.AccountState)

	account := account.NewAddress()
	err = account.UnmarshalBinary(state.Account)
	if err != nil {
		return err
	}

	s.StateType = libblock.StateType(core.CORE_ACCOUNT_STATE)
	s.BlockIndex = state.BlockIndex
	s.Account = account
	s.Sequence = state.Sequence
	s.Amount = state.Amount
	return nil
}

func (s *AccountState) MarshalBinary() ([]byte, error) {
	a, err := s.Account.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return core.Marshal(&pb.AccountState{
		StateType:  uint32(core.CORE_ACCOUNT_STATE),
		BlockIndex: s.BlockIndex,
		Account:    a,
		Sequence:   s.Sequence,
		Amount:     s.Amount,
	})
}

func (s *AccountState) Raw(ignoreSigningFields bool) ([]byte, error) {
	a, err := s.Account.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return core.Marshal(&pb.AccountState{
		StateType: uint32(core.CORE_ACCOUNT_STATE),
		Account:   a,
		Sequence:  s.Sequence,
		Amount:    s.Amount,
	})
}

type CurrencyState struct {
	State

	Account  libcore.Address
	Sequence uint64

	Name        string
	Symbol      string
	Decimals    uint32
	TotalSupply int64
}

func (s *CurrencyState) GetIndex() uint64 {
	return s.Sequence
}

func (s *CurrencyState) GetStateKey() string {
	return s.Symbol
}

func (s *CurrencyState) UnmarshalBinary(data []byte) error {
	meta, msg, err := core.Unmarshal(data)
	if err != nil {
		return err
	}
	if meta != core.CORE_CURRENCY_STATE {
		return errors.New("error state data")
	}
	state := msg.(*pb.CurrencyState)

	issuer := account.NewAddress()
	err = issuer.UnmarshalBinary(state.Account)
	if err != nil {
		return err
	}

	s.StateType = libblock.StateType(core.CORE_CURRENCY_STATE)
	s.BlockIndex = state.BlockIndex
	s.Account = issuer
	s.Sequence = state.Sequence
	s.Name = state.Name
	s.Symbol = state.Symbol
	s.Decimals = state.Decimals
	s.TotalSupply = state.TotalSupply

	return nil
}

func (s *CurrencyState) MarshalBinary() ([]byte, error) {
	issuer, err := s.Account.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return core.Marshal(&pb.CurrencyState{
		StateType:   uint32(core.CORE_CURRENCY_STATE),
		BlockIndex:  s.BlockIndex,
		Account:     issuer,
		Sequence:    s.Sequence,
		Name:        s.Name,
		Symbol:      s.Symbol,
		Decimals:    s.Decimals,
		TotalSupply: s.TotalSupply,
	})
}

func (s *CurrencyState) Raw(ignoreSigningFields bool) ([]byte, error) {
	issuer, err := s.Account.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return core.Marshal(&pb.CurrencyState{
		StateType:   uint32(core.CORE_CURRENCY_STATE),
		Account:     issuer,
		Sequence:    s.Sequence,
		Name:        s.Name,
		Symbol:      s.Symbol,
		Decimals:    s.Decimals,
		TotalSupply: s.TotalSupply,
	})
}

type DeviceState struct {
	State

	Account  libcore.Address
	Sequence uint64

	Symbol      string
	Description string
	Tags        []string
}

func (s *DeviceState) GetIndex() uint64 {
	return s.Sequence
}

func (s *DeviceState) GetStateKey() string {
	return s.Symbol
}

func (s *DeviceState) UnmarshalBinary(data []byte) error {
	meta, msg, err := core.Unmarshal(data)
	if err != nil {
		return err
	}
	if meta != core.CORE_DEVICE_STATE {
		return errors.New("error state data")
	}
	state := msg.(*pb.DeviceState)

	issuer := account.NewAddress()
	err = issuer.UnmarshalBinary(state.Account)
	if err != nil {
		return err
	}

	s.StateType = libblock.StateType(core.CORE_DEVICE_STATE)
	s.BlockIndex = state.BlockIndex
	s.Account = issuer
	s.Sequence = state.Sequence
	s.Symbol = state.Symbol
	s.Description = state.Description
	s.Tags = state.Tags

	return nil
}

func (s *DeviceState) MarshalBinary() ([]byte, error) {
	issuer, err := s.Account.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return core.Marshal(&pb.DeviceState{
		StateType:   uint32(core.CORE_DEVICE_STATE),
		BlockIndex:  s.BlockIndex,
		Account:     issuer,
		Sequence:    s.Sequence,
		Symbol:      s.Symbol,
		Description: s.Description,
		Tags:        s.Tags,
	})
}

func (s *DeviceState) Raw(ignoreSigningFields bool) ([]byte, error) {
	issuer, err := s.Account.MarshalBinary()
	if err != nil {
		return nil, err
	}

	return core.Marshal(&pb.DeviceState{
		StateType:   uint32(core.CORE_DEVICE_STATE),
		Account:     issuer,
		Sequence:    s.Sequence,
		Symbol:      s.Symbol,
		Description: s.Description,
		Tags:        s.Tags,
	})
}

func ReadState(data []byte) (libblock.State, error) {
	if len(data) == 0 {
		return nil, errors.New("error entry")
	}
	meta := data[0]
	switch meta {
	case core.CORE_ACCOUNT_STATE:
		s := &AccountState{}
		err := s.UnmarshalBinary(data)
		if err != nil {
			return nil, err
		}
		return s, nil
	case core.CORE_CURRENCY_STATE:
		s := &CurrencyState{}
		err := s.UnmarshalBinary(data)
		if err != nil {
			return nil, err
		}
		return s, nil
	case core.CORE_DEVICE_STATE:
		s := &DeviceState{}
		err := s.UnmarshalBinary(data)
		if err != nil {
			return nil, err
		}
		return s, nil
	default:
		return nil, errors.New("error data")
	}
}

func CloneState(state libblock.State) (libblock.State, error) {
	data, err := state.MarshalBinary()
	if err != nil {
		return nil, err
	}
	s, err := ReadState(data)
	if err != nil {
		return nil, err
	}
	return s, nil
}
