package kinetic

import (
	kproto "github.com/yongzhy/kinetic-go/proto"
)

// Callback is the interface define actions for MessageType.
// Success is called when XXXXX_RESPONSE message recieved from drive without problem.
// Failure is called when XXXXX_RESPONSE message status code is not OK, or any other kind of failure.
// Done return true if either Success or Failure is called to indicate XXXXX_RESPONSE received and processed.
// Status return the MessateType operation status.
type Callback interface {
	Success(resp *kproto.Command, value []byte)
	Failure(status Status)
	Status() Status
}

// GenericCallback can be used for all MessageType which doesn't require data from Kinetic drive.
// And for MessageType that require data from drive, a new struct need to be defined GenericCallback
type GenericCallback struct {
	status Status
}

// Success is called by ResponseHandler when response message received from kinetic device has OK status.
func (c *GenericCallback) Success(resp *kproto.Command, value []byte) {
	c.status = Status{Code: OK}
}

// Failure is called ResponseHandler when response message received from kinetic device with status code other than OK.
func (c *GenericCallback) Failure(status Status) {
	c.status = status
}

// Status returns the status after ResponseHandler processed response message from kinetic device.
func (c *GenericCallback) Status() Status {
	return c.status
}

// GetCallback is the Callback for Command_GET Message
type GetCallback struct {
	GenericCallback
	Entry Record // Entity information
}

// Success function extracts object information from response message and
// store into GetCallback.Entry.
func (c *GetCallback) Success(resp *kproto.Command, value []byte) {
	c.GenericCallback.Success(resp, value)
	c.Entry.Key = resp.GetBody().GetKeyValue().GetKey()
	c.Entry.Tag = resp.GetBody().GetKeyValue().GetTag()
	c.Entry.Version = resp.GetBody().GetKeyValue().GetDbVersion()
	c.Entry.Algo = convertAlgoFromProto(resp.GetBody().GetKeyValue().GetAlgorithm())

	c.Entry.Value = value
}

// GetKeyRangeCallback is the Callback for Command_GETKEYRANGE Message
type GetKeyRangeCallback struct {
	GenericCallback
	Keys [][]byte // List of objects' keys within range, get from device
}

// Success extracts objects' keys within range, from response message.
func (c *GetKeyRangeCallback) Success(resp *kproto.Command, value []byte) {
	c.GenericCallback.Success(resp, value)
	c.Keys = resp.GetBody().GetRange().GetKeys()
}

// GetVersionCallback is the Callback for Command_GETVERSION Message
type GetVersionCallback struct {
	GenericCallback
	Version []byte // Version of the object on device
}

// Success extracts object's version information from response message.
func (c *GetVersionCallback) Success(resp *kproto.Command, value []byte) {
	c.GenericCallback.Success(resp, value)
	c.Version = resp.GetBody().GetKeyValue().GetDbVersion()
}

// P2PPushCallback is the Callback for Command_PEER2PEERPUSH
type P2PPushCallback struct {
	GenericCallback
	Statuses []Status
}

// Success extracts P2Push operation status from response message.
func (c *P2PPushCallback) Success(resp *kproto.Command, value []byte) {
	c.GenericCallback.Success(resp, value)
	c.Statuses = make([]Status, len(resp.GetBody().GetP2POperation().GetOperation()))
	for k, op := range resp.GetBody().GetP2POperation().GetOperation() {
		c.Statuses[k].Code = convertStatusCodeFromProto(op.GetStatus().GetCode())
		c.Statuses[k].ErrorMsg = op.GetStatus().GetStatusMessage()
	}
}

// GetLogCallback is the Callback for Command_GETLOG Message
type GetLogCallback struct {
	GenericCallback
	Logs Log // Device log information
}

// Success extracts kientic device's Log information from response message.
func (c *GetLogCallback) Success(resp *kproto.Command, value []byte) {
	c.GenericCallback.Success(resp, value)
	c.Logs = getLogFromProto(resp)
}
