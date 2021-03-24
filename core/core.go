package core

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/tokentransfer/chain/core/pb"
	"google.golang.org/protobuf/proto"
)

const (
	CORE_BLOCK                 = byte(100)
	CORE_TRANSACTION           = byte(101)
	CORE_RECEIPT               = byte(102)
	CORE_TRANSACTION_WITH_DATA = byte(103)
	CORE_PAYMENT               = byte(104)
	CORE_NEWDEVICE             = byte(105)
	CORE_PAYMENT_WITH_DATA     = byte(106)
	CORE_NEWDEVICE_WITH_DATA   = byte(107)

	// CORE_STATE         = byte(110)
	CORE_ACCOUNT_STATE  = byte(111)
	CORE_CURRENCY_STATE = byte(112)
	CORE_DEVICE_STATE   = byte(113)

	CORE_PAYMENT_TYPE      = byte(201)
	CORE_NEW_CURRENCY_TYPE = byte(202)
	CORE_NEW_DEVICE_TYPE   = byte(203)
)

func GetInfo(data []byte) string {
	if len(data) > 0 {
		meta := data[0]
		switch meta {
		case CORE_BLOCK:
			return "block"
		case CORE_TRANSACTION:
			return "transaction"
		case CORE_PAYMENT:
			return "payment"
		case CORE_NEWDEVICE:
			return "new_device"
		case CORE_RECEIPT:
			return "receipt"
		case CORE_TRANSACTION_WITH_DATA:
			return "transaction_with_data"
		case CORE_PAYMENT_WITH_DATA:
			return "payment_with_data"
		case CORE_NEWDEVICE_WITH_DATA:
			return "newdevice_with_data"

		case CORE_ACCOUNT_STATE:
			return "account_state"
		case CORE_CURRENCY_STATE:
			return "currency_state"
		case CORE_DEVICE_STATE:
			return "device_state"

		case CORE_PAYMENT_TYPE:
			return "payment_type"
		case CORE_NEW_CURRENCY_TYPE:
			return "new_currency_type"
		case CORE_NEW_DEVICE_TYPE:
			return "new_device_type"

		default:
			return "unknown"
		}
	}
	return ""
}

func Clone(t proto.Message) (proto.Message, error) {
	data, err := Marshal(t)
	if err != nil {
		return nil, err
	}
	_, o, err := Unmarshal(data)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func Marshal(message proto.Message) ([]byte, error) {
	var meta byte
	switch message.(type) {
	case *pb.Block:
		meta = CORE_BLOCK
	case *pb.Transaction:
		meta = CORE_TRANSACTION
	case *pb.Payment:
		meta = CORE_PAYMENT
	case *pb.NewDevice:
		meta = CORE_NEWDEVICE
	case *pb.Receipt:
		meta = CORE_RECEIPT
	case *pb.TransactionWithData:
		meta = CORE_TRANSACTION_WITH_DATA
	case *pb.PaymentWithData:
		meta = CORE_PAYMENT_WITH_DATA
	case *pb.NewDeviceWithData:
		meta = CORE_NEWDEVICE_WITH_DATA

	case *pb.AccountState:
		meta = CORE_ACCOUNT_STATE
	case *pb.CurrencyState:
		meta = CORE_CURRENCY_STATE
	case *pb.DeviceState:
		meta = CORE_DEVICE_STATE

	default:
		err := errors.New("error data type")
		return nil, err
	}
	data, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}
	return append([]byte{meta}, data...), nil
}

func Unmarshal(data []byte) (byte, proto.Message, error) {
	meta := data[0]
	bs := data[1:]

	var msg proto.Message
	switch meta {
	case CORE_BLOCK:
		msg = &pb.Block{}
	case CORE_TRANSACTION:
		msg = &pb.Transaction{}
	case CORE_PAYMENT:
		msg = &pb.Payment{}
	case CORE_NEWDEVICE:
		msg = &pb.NewDevice{}
	case CORE_RECEIPT:
		msg = &pb.Receipt{}
	case CORE_TRANSACTION_WITH_DATA:
		msg = &pb.TransactionWithData{}
	case CORE_PAYMENT_WITH_DATA:
		msg = &pb.PaymentWithData{}
	case CORE_NEWDEVICE_WITH_DATA:
		msg = &pb.NewDeviceWithData{}

	case CORE_ACCOUNT_STATE:
		msg = &pb.AccountState{}
	case CORE_CURRENCY_STATE:
		msg = &pb.CurrencyState{}
	case CORE_DEVICE_STATE:
		msg = &pb.DeviceState{}

	default:
		err := errors.New("error data format")
		return 0, nil, err
	}

	err := proto.Unmarshal(bs, msg)
	if err != nil {
		return 0, nil, err
	}
	return meta, msg, nil
}

func WriteBytes(w io.Writer, b []byte) error {
	l := len(b)
	err := binary.Write(w, binary.LittleEndian, uint32(l))
	if err != nil {
		return err
	}
	if l > 0 {
		n, err := w.Write(b)
		if err != nil {
			return err
		}
		if n != l {
			return errors.New("error write")
		}
	}
	return nil
}

func ReadBytes(r io.Reader) ([]byte, error) {
	l := uint32(0)
	err := binary.Read(r, binary.LittleEndian, &l)
	if err != nil {
		return nil, err
	}
	b := make([]byte, l)
	if l > 0 {
		n, err := r.Read(b)
		if err != nil {
			return nil, err
		}
		if n != int(l) {
			return nil, errors.New("error read")
		}
	}
	return b, nil
}
