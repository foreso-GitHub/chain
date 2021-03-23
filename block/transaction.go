package block

import (
	"errors"
	"log"

	"github.com/tokentransfer/chain/account"
	"github.com/tokentransfer/chain/core"
	"github.com/tokentransfer/chain/core/pb"

	libblock "github.com/tokentransfer/interfaces/block"
	libcore "github.com/tokentransfer/interfaces/core"
)

type Transaction struct {
	Hash libcore.Hash

	TransactionType libblock.TransactionType

	Account  libcore.Address
	Sequence uint64
	Amount   int64
	Gas      int64
	Type     string

	//Timestamp   int64
	//Tags        []string
	//Name        string
	//Value       string
	//Device      string

	//Symbol      string
	//Description string
	//DeviceTags  []string

	Destination libcore.Address
	Payload     libcore.Bytes
	PublicKey   libcore.PublicKey
	Signature   libcore.Signature
}

func (tx *Transaction) GetIndex() uint64 {
	return tx.Sequence
}

func (tx *Transaction) GetHash() libcore.Hash {
	return tx.Hash
}

func (tx *Transaction) SetHash(h libcore.Hash) {
	tx.Hash = h
}

func byteToAddress(b []byte) (libcore.Address, error) {
	a := account.NewAddress()
	err := a.UnmarshalBinary(b)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (tx *Transaction) UnmarshalBinary(data []byte) error {
	var err error

	meta, msg, err := core.Unmarshal(data)
	if err != nil {
		return err
	}
	if meta != core.CORE_TRANSACTION {
		return errors.New("error transaction data")
	}
	t := msg.(*pb.Transaction)

	tx.TransactionType = libblock.TransactionType(t.TransactionType)

	tx.Account, err = byteToAddress(t.Account)
	if err != nil {
		return err
	}

	tx.Sequence = t.Sequence
	tx.Amount = t.Amount
	tx.Gas = t.Gas
	tx.Type = t.Type

	tx.Destination, err = byteToAddress(t.Destination)
	if err != nil {
		return err
	}

	tx.Payload = t.Payload
	tx.PublicKey = libcore.PublicKey(t.PublicKey)
	tx.Signature = libcore.Signature(t.Signature)

	return nil
}

func addressToByte(a libcore.Address) ([]byte, error) {
	return a.MarshalBinary()
}

func (tx *Transaction) MarshalBinary() ([]byte, error) {
	fromData, err := addressToByte(tx.Account)
	if err != nil {
		return nil, err
	}
	toData, err := addressToByte(tx.Destination)
	if err != nil {
		return nil, err
	}

	t := &pb.Transaction{
		TransactionType: uint32(tx.TransactionType),

		Account:     fromData,
		Sequence:    tx.Sequence,
		Amount:      tx.Amount,
		Gas:         tx.Gas,
		Type:        tx.Type,
		Destination: toData,
		Payload:     tx.Payload,
		PublicKey:   []byte(tx.PublicKey),
		Signature:   []byte(tx.Signature),
	}
	return core.Marshal(t)
}

func (tx *Transaction) Raw(ignoreSigningFields bool) ([]byte, error) {
	fromData, err := addressToByte(tx.Account)
	if err != nil {
		return nil, err
	}
	toData, err := addressToByte(tx.Destination)
	if err != nil {
		return nil, err
	}

	if ignoreSigningFields {
		t := &pb.Transaction{
			TransactionType: uint32(tx.TransactionType),

			Account:     fromData,
			Sequence:    tx.Sequence,
			Amount:      tx.Amount,
			Gas:         tx.Gas,
			Type:        tx.Type,
			Destination: toData,
			Payload:     tx.Payload,
			PublicKey:   []byte(tx.PublicKey),
		}
		return core.Marshal(t)
	}
	return tx.MarshalBinary()
}

func (tx *Transaction) GetTransactionType() libblock.TransactionType {
	return tx.TransactionType
}

func (tx *Transaction) GetAccount() libcore.Address {
	return tx.Account
}

func (tx *Transaction) GetPublicKey() libcore.PublicKey {
	return tx.PublicKey
}

func (tx *Transaction) SetPublicKey(p libcore.PublicKey) {
	tx.PublicKey = p
}

func (tx *Transaction) GetSignature() libcore.Signature {
	return tx.Signature
}

func (tx *Transaction) SetSignature(s libcore.Signature) {
	tx.Signature = s
}

//region BaseTx

//type BaseTx interface {
//	libblock.Transaction
//
//	getAccount()		libcore.Address
//	getDistination()	libcore.Address
//	getSequence()		uint64
//	getAmount()			int64
//	getGas()			int64
//	getType()			string
//}

//func (tx *Transaction) GetAccount() libcore.Address {
//	return tx.Account
//}

func (tx *Transaction) GetDestination() libcore.Address {
	return tx.Destination
}

func (tx *Transaction) GetSequence() uint64 {
	return tx.Sequence
}

func (tx *Transaction) GetAmount() int64 {
	return tx.Amount
}

func (tx *Transaction) GetGas() int64 {
	return tx.Gas
}

func (tx *Transaction) GetType() string {
	return tx.Type
}

//endregion

//region TransactionWithData

type TransactionWithData struct {
	Hash libcore.Hash

	Transaction libblock.Transaction
	Receipt     libblock.Receipt
	Payment     libblock.Transaction
	NewDevice   libblock.Transaction
}

func (txWithData *TransactionWithData) GetHash() libcore.Hash {
	return txWithData.Hash
}

func (txWithData *TransactionWithData) SetHash(h libcore.Hash) {
	txWithData.Hash = h
}

func (txWithData *TransactionWithData) GetTransaction() libblock.Transaction {
	return txWithData.Transaction
}

func (txWithData *TransactionWithData) GetReceipt() libblock.Receipt {
	return txWithData.Receipt
}

func (txWithData *TransactionWithData) UnmarshalBinary(data []byte) error {
	meta, msg, err := core.Unmarshal(data)
	if meta != core.CORE_TRANSACTION_WITH_DATA {
		return errors.New("error transaction with data")
	}

	td := msg.(*pb.TransactionWithData)

	txData, err := core.Marshal(td.Transaction)
	if err != nil {
		return err
	}
	tx := &Transaction{}
	err = tx.UnmarshalBinary(txData)
	if err == nil {
		txWithData.Transaction = tx
	}

	paymentData, err := core.Marshal(td.Payment)
	if err != nil {
		return err
	}
	payment := &Payment{}
	err = tx.UnmarshalBinary(paymentData)
	if err == nil {
		txWithData.Payment = payment
	}

	deviceData, err := core.Marshal(td.NewDevice)
	if err != nil {
		return err
	}
	device := &NewDevice{}
	err = tx.UnmarshalBinary(deviceData)
	if err == nil {
		txWithData.NewDevice = device
	}

	receiptData, err := core.Marshal(td.Receipt)
	if err != nil {
		log.Println(err)
		return err
	}
	receipt := &Receipt{}
	err = receipt.UnmarshalBinary(receiptData)
	if err != nil {
		log.Println(err)
		return err
	}

	txWithData.Receipt = receipt
	return nil
}

func (txWithData *TransactionWithData) MarshalBinary() ([]byte, error) {

	receiptData, err := txWithData.Receipt.MarshalBinary()
	if err != nil {
		return nil, err
	}
	_, msg, err := core.Unmarshal(receiptData)
	if err != nil {
		return nil, err
	}
	receipt := msg.(*pb.Receipt)

	txData, err := txWithData.Transaction.MarshalBinary()
	if err != nil {
		return nil, err
	}
	meta, msg, err := core.Unmarshal(txData)
	if err != nil {
		return nil, err
	}

	//tx := msg.(*pb.Transaction)
	//tx := msg.(*pb.Payment)

	switch meta {
	case core.CORE_TRANSACTION:
		tx := msg.(*pb.Transaction)

		td := &pb.TransactionWithData{
			Transaction: tx,
			Receipt:     receipt,
		}
		data, err := core.Marshal(td)
		if err != nil {
			return nil, err
		}
		return data, nil
	case core.CORE_PAYMENT:
		payment := msg.(*pb.Payment)
		td := &pb.TransactionWithData{
			Payment: payment,
			Receipt: receipt,
		}
		data, err := core.Marshal(td)
		if err != nil {
			return nil, err
		}
		return data, nil
	case core.CORE_NEWDEVICE:
		newDevice := msg.(*pb.NewDevice)
		td := &pb.TransactionWithData{
			NewDevice: newDevice,
			Receipt:   receipt,
		}
		data, err := core.Marshal(td)
		if err != nil {
			return nil, err
		}
		return data, nil
	default:
		err := errors.New("error TransactionWithData meta")
		return nil, err
	}

	//td := &pb.TransactionWithData{
	//	Transaction: tx,
	//	Receipt:     receipt,
	//}
	//data, err := core.Marshal(td)
	//if err != nil {
	//	return nil, err
	//}
	//return data, nil
}

func (txWithData *TransactionWithData) Raw(ignoreSigningFields bool) ([]byte, error) {

	receiptData, err := txWithData.Receipt.Raw(ignoreSigningFields)
	if err != nil {
		return nil, err
	}
	_, msg, err := core.Unmarshal(receiptData)
	if err != nil {
		return nil, err
	}
	receipt := msg.(*pb.Receipt)

	txData, err := txWithData.Transaction.Raw(ignoreSigningFields)
	if err != nil {
		return nil, err
	}
	meta, msg, err := core.Unmarshal(txData)
	if err != nil {
		return nil, err
	}
	//tx := msg.(*pb.Transaction)
	//tx := msg.(*pb.Payment)

	switch meta {
	case core.CORE_TRANSACTION:
		tx := msg.(*pb.Transaction)

		td := &pb.TransactionWithData{
			Transaction: tx,
			Receipt:     receipt,
		}
		data, err := core.Marshal(td)
		if err != nil {
			return nil, err
		}
		return data, nil
	case core.CORE_PAYMENT:
		payment := msg.(*pb.Payment)
		td := &pb.TransactionWithData{
			Payment: payment,
			Receipt: receipt,
		}
		data, err := core.Marshal(td)
		if err != nil {
			return nil, err
		}
		return data, nil
	case core.CORE_NEWDEVICE:
		newDevice := msg.(*pb.NewDevice)
		td := &pb.TransactionWithData{
			NewDevice: newDevice,
			Receipt:   receipt,
		}
		data, err := core.Marshal(td)
		if err != nil {
			return nil, err
		}
		return data, nil
	default:
		err := errors.New("error TransactionWithData meta")
		return nil, err
	}

	//td := &pb.TransactionWithData{
	//	Transaction: tx,
	//	Receipt:     receipt,
	//}
	//data, err := core.Marshal(td)
	//if err != nil {
	//	return nil, err
	//}
	//return data, nil
}

//endregion

//region Payment

type Payment struct {
	Transaction

	Timestamp int64
	Device    string
	Tags      []string
	Name      string
	Value     string
}

func (tx *Payment) UnmarshalBinary(data []byte) error {
	var err error

	meta, msg, err := core.Unmarshal(data)
	if err != nil {
		return err
	}
	if meta != core.CORE_PAYMENT {
		return errors.New("error transaction payment data")
	}
	t := msg.(*pb.Payment)

	tx.TransactionType = libblock.TransactionType(t.TransactionType)

	tx.Account, err = byteToAddress(t.Account)
	if err != nil {
		return err
	}

	tx.Sequence = t.Sequence
	tx.Amount = t.Amount
	tx.Gas = t.Gas
	tx.Type = t.Type
	tx.Timestamp = t.Timestamp
	tx.Tags = t.Tags
	tx.Name = t.Name
	tx.Value = t.Value
	tx.Device = t.Device

	tx.Destination, err = byteToAddress(t.Destination)
	if err != nil {
		return err
	}

	tx.Payload = t.Payload
	tx.PublicKey = libcore.PublicKey(t.PublicKey)
	tx.Signature = libcore.Signature(t.Signature)

	return nil
}

func (tx *Payment) MarshalBinary() ([]byte, error) {
	fromData, err := addressToByte(tx.Account)
	if err != nil {
		return nil, err
	}
	toData, err := addressToByte(tx.Destination)
	if err != nil {
		return nil, err
	}

	t := &pb.Payment{
		TransactionType: uint32(tx.TransactionType),

		Account:     fromData,
		Sequence:    tx.Sequence,
		Amount:      tx.Amount,
		Gas:         tx.Gas,
		Timestamp:   tx.Timestamp,
		Tags:        tx.Tags,
		Name:        tx.Name,
		Value:       tx.Value,
		Device:      tx.Device,
		Type:        tx.Type,
		Destination: toData,
		Payload:     tx.Payload,
		PublicKey:   []byte(tx.PublicKey),
		Signature:   []byte(tx.Signature),
	}
	return core.Marshal(t)
}

func (tx *Payment) Raw(ignoreSigningFields bool) ([]byte, error) {
	fromData, err := addressToByte(tx.Account)
	if err != nil {
		return nil, err
	}
	toData, err := addressToByte(tx.Destination)
	if err != nil {
		return nil, err
	}

	if ignoreSigningFields {
		t := &pb.Payment{
			TransactionType: uint32(tx.TransactionType),

			Account:     fromData,
			Sequence:    tx.Sequence,
			Amount:      tx.Amount,
			Gas:         tx.Gas,
			Timestamp:   tx.Timestamp,
			Tags:        tx.Tags,
			Name:        tx.Name,
			Value:       tx.Value,
			Device:      tx.Device,
			Type:        tx.Type,
			Destination: toData,
			Payload:     tx.Payload,
			PublicKey:   []byte(tx.PublicKey),
		}
		return core.Marshal(t)
	}
	return tx.MarshalBinary()
}

//endregion

//region NewDevice

type NewDevice struct {
	Transaction

	Symbol      string
	Description string
	DeviceTags  []string
}

func (tx *NewDevice) UnmarshalBinary(data []byte) error {
	var err error

	meta, msg, err := core.Unmarshal(data)
	if err != nil {
		return err
	}
	if meta != core.CORE_NEWDEVICE {
		return errors.New("error transaction new device data")
	}
	t := msg.(*pb.NewDevice)

	tx.TransactionType = libblock.TransactionType(t.TransactionType)

	tx.Account, err = byteToAddress(t.Account)
	if err != nil {
		return err
	}

	tx.Sequence = t.Sequence
	tx.Amount = t.Amount
	tx.Gas = t.Gas
	tx.Type = t.Type
	tx.Symbol = t.Symbol
	tx.Description = t.Description
	tx.DeviceTags = t.DeviceTags

	tx.Destination, err = byteToAddress(t.Destination)
	if err != nil {
		return err
	}

	tx.Payload = t.Payload
	tx.PublicKey = libcore.PublicKey(t.PublicKey)
	tx.Signature = libcore.Signature(t.Signature)

	return nil
}

func (tx *NewDevice) MarshalBinary() ([]byte, error) {
	fromData, err := addressToByte(tx.Account)
	if err != nil {
		return nil, err
	}
	toData, err := addressToByte(tx.Destination)
	if err != nil {
		return nil, err
	}

	t := &pb.NewDevice{
		TransactionType: uint32(tx.TransactionType),

		Account:     fromData,
		Sequence:    tx.Sequence,
		Amount:      tx.Amount,
		Gas:         tx.Gas,
		Type:        tx.Type,
		Symbol:      tx.Symbol,
		Description: tx.Description,
		DeviceTags:  tx.DeviceTags,
		Destination: toData,
		Payload:     tx.Payload,
		PublicKey:   []byte(tx.PublicKey),
		Signature:   []byte(tx.Signature),
	}
	return core.Marshal(t)
}

func (tx *NewDevice) Raw(ignoreSigningFields bool) ([]byte, error) {
	fromData, err := addressToByte(tx.Account)
	if err != nil {
		return nil, err
	}
	toData, err := addressToByte(tx.Destination)
	if err != nil {
		return nil, err
	}

	if ignoreSigningFields {
		t := &pb.NewDevice{
			TransactionType: uint32(tx.TransactionType),

			Account:     fromData,
			Sequence:    tx.Sequence,
			Amount:      tx.Amount,
			Gas:         tx.Gas,
			Type:        tx.Type,
			Symbol:      tx.Symbol,
			Description: tx.Description,
			DeviceTags:  tx.DeviceTags,
			Destination: toData,
			Payload:     tx.Payload,
			PublicKey:   []byte(tx.PublicKey),
		}
		return core.Marshal(t)
	}
	return tx.MarshalBinary()
}

//endregion
