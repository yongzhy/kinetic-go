package kinetic

import (
	kproto "github.com/yongzhy/kinetic-go/proto"
)

type NonBlockConnection struct {
	service *networkService
}

func NewNonBlockConnection(op ClientOptions) (*NonBlockConnection, error) {
	if op.Hmac == nil {
		klog.Panic("HMAC is required for ClientOptions")
	}

	service, err := newNetworkService(op)
	if err != nil {
		return nil, err
	}

	return &NonBlockConnection{service}, nil
}

func (conn *NonBlockConnection) NoOp(h *MessageHandler) error {
	msg := newMessage(kproto.Message_HMACAUTH)

	cmd := newCommand(kproto.Command_NOOP)

	return conn.service.submit(msg, cmd, nil, h)
}

func (conn *NonBlockConnection) get(key []byte, getType kproto.Command_MessageType, h *MessageHandler) error {
	msg := newMessage(kproto.Message_HMACAUTH)

	cmd := newCommand(getType)
	cmd.Body = &kproto.Command_Body{
		KeyValue: &kproto.Command_KeyValue{
			Key: key,
		},
	}

	return conn.service.submit(msg, cmd, nil, h)
}

func (conn *NonBlockConnection) Get(key []byte, h *MessageHandler) error {
	return conn.get(key, kproto.Command_GET, h)
}

func (conn *NonBlockConnection) GetNext(key []byte, h *MessageHandler) error {
	return conn.get(key, kproto.Command_GETNEXT, h)
}

func (conn *NonBlockConnection) GetPrevious(key []byte, h *MessageHandler) error {
	return conn.get(key, kproto.Command_GETPREVIOUS, h)
}

func (conn *NonBlockConnection) GetKeyRange(r *KeyRange, h *MessageHandler) error {
	msg := newMessage(kproto.Message_HMACAUTH)

	cmd := newCommand(kproto.Command_GETKEYRANGE)
	cmd.Body = &kproto.Command_Body{
		Range: &kproto.Command_Range{
			StartKey:          r.StartKey,
			EndKey:            r.EndKey,
			StartKeyInclusive: &r.StartKeyInclusive,
			EndKeyInclusive:   &r.EndKeyInclusive,
			MaxReturned:       &r.Max,
			Reverse:           &r.Reverse,
		},
	}

	return conn.service.submit(msg, cmd, nil, h)
}

func (conn *NonBlockConnection) Delete(entry *Record, h *MessageHandler) error {
	msg := newMessage(kproto.Message_HMACAUTH)
	cmd := newCommand(kproto.Command_DELETE)

	sync := convertSyncToProto(entry.Sync)
	//algo := convertAlgoToProto(entry.Algo)
	cmd.Body = &kproto.Command_Body{
		KeyValue: &kproto.Command_KeyValue{
			Key:             entry.Key,
			Force:           &entry.Force,
			Synchronization: &sync,
			//Algorithm:       &algo,
		},
	}

	return conn.service.submit(msg, cmd, nil, h)
}

func (conn *NonBlockConnection) Put(entry *Record, h *MessageHandler) error {
	msg := newMessage(kproto.Message_HMACAUTH)
	cmd := newCommand(kproto.Command_PUT)

	sync := convertSyncToProto(entry.Sync)
	algo := convertAlgoToProto(entry.Algo)
	cmd.Body = &kproto.Command_Body{
		KeyValue: &kproto.Command_KeyValue{
			Key:             entry.Key,
			Force:           &entry.Force,
			Synchronization: &sync,
			Algorithm:       &algo,
			Tag:             entry.Tag,
		},
	}

	return conn.service.submit(msg, cmd, entry.Value, h)
}

func (conn *NonBlockConnection) pinop(pin []byte, op kproto.Command_PinOperation_PinOpType, h *MessageHandler) error {
	msg := newMessage(kproto.Message_PINAUTH)
	msg.PinAuth = &kproto.Message_PINauth{
		Pin: pin,
	}

	cmd := newCommand(kproto.Command_PINOP)

	cmd.Body = &kproto.Command_Body{
		PinOp: &kproto.Command_PinOperation{
			PinOpType: &op,
		},
	}

	return conn.service.submit(msg, cmd, nil, h)
}

func (conn *NonBlockConnection) SecureErase(pin []byte, h *MessageHandler) error {
	return conn.pinop(pin, kproto.Command_PinOperation_SECURE_ERASE_PINOP, h)
}

func (conn *NonBlockConnection) InstantErase(pin []byte, h *MessageHandler) error {
	return conn.pinop(pin, kproto.Command_PinOperation_ERASE_PINOP, h)

}

func (conn *NonBlockConnection) LockDevice(pin []byte, h *MessageHandler) error {
	return conn.pinop(pin, kproto.Command_PinOperation_LOCK_PINOP, h)
}

func (conn *NonBlockConnection) UnlockDevice(pin []byte, h *MessageHandler) error {
	return conn.pinop(pin, kproto.Command_PinOperation_UNLOCK_PINOP, h)
}

func (conn *NonBlockConnection) UpdateFirmware(code []byte, h *MessageHandler) error {
	msg := newMessage(kproto.Message_HMACAUTH)
	cmd := newCommand(kproto.Command_SETUP)

	var download bool = true
	cmd.Body = &kproto.Command_Body{
		Setup: &kproto.Command_Setup{
			FirmwareDownload: &download,
		},
	}

	return conn.service.submit(msg, cmd, code, h)
}

func (conn *NonBlockConnection) SetClusterVersion(version int64, h *MessageHandler) error {
	msg := newMessage(kproto.Message_HMACAUTH)
	cmd := newCommand(kproto.Command_SETUP)

	cmd.Body = &kproto.Command_Body{
		Setup: &kproto.Command_Setup{
			NewClusterVersion: &version,
		},
	}

	return conn.service.submit(msg, cmd, nil, h)
}

func (conn *NonBlockConnection) SetLockPin(currentPin []byte, newPin []byte, h *MessageHandler) error {
	msg := newMessage(kproto.Message_HMACAUTH)
	cmd := newCommand(kproto.Command_SECURITY)

	cmd.Body = &kproto.Command_Body{
		Security: &kproto.Command_Security{
			OldLockPIN: currentPin,
			NewLockPIN: newPin,
		},
	}

	return conn.service.submit(msg, cmd, nil, h)
}

func (conn *NonBlockConnection) SetErasePin(currentPin []byte, newPin []byte, h *MessageHandler) error {
	msg := newMessage(kproto.Message_HMACAUTH)
	cmd := newCommand(kproto.Command_SECURITY)

	cmd.Body = &kproto.Command_Body{
		Security: &kproto.Command_Security{
			OldErasePIN: currentPin,
			NewErasePIN: newPin,
		},
	}

	return conn.service.submit(msg, cmd, nil, h)
}

func (conn *NonBlockConnection) Run() error {
	return conn.service.listen()
}

func (conn *NonBlockConnection) Close() {
	conn.service.close()
}
