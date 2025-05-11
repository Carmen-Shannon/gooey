//go:build windows
// +build windows

package wdws

import "golang.org/x/sys/windows"

type getMessageOpts struct {
	WindowHandle     windows.Handle
	Message          *Msg
	MessageFilterMin uint32
	MessageFilterMax uint32
	DoneChan         chan struct{}
}

type GetMessageOpt func(*getMessageOpts)

func newGetMessageOpts() *getMessageOpts {
	return &getMessageOpts{
		WindowHandle:     0,
		Message:          nil,
		MessageFilterMin: 0,
		MessageFilterMax: 0,
	}
}

// WindowHandleOpt sets the window handle for the get message options.
//
// Parameters:
//   - windowHandle: The handle to the window to receive messages.
//
// Returns:
//   - GetMessageOpt: A function that takes a pointer to getMessageOpts and sets its WindowHandle field.
func WindowHandleOpt(windowHandle windows.Handle) GetMessageOpt {
    return func(opts *getMessageOpts) {
        opts.WindowHandle = windowHandle
    }
}

// MessageOpt sets the message pointer for the get message options.
//
// Parameters:
//   - message: The pointer to the Msg structure to receive the message.
//
// Returns:
//   - GetMessageOpt: A function that takes a pointer to getMessageOpts and sets its Message field.
func MessageOpt(message *Msg) GetMessageOpt {
    return func(opts *getMessageOpts) {
        opts.Message = message
    }
}

// MessageFilterMinOpt sets the minimum message filter value for the get message options.
//
// Parameters:
//   - messageFilterMin: The minimum value of the message filter range.
//
// Returns:
//   - GetMessageOpt: A function that takes a pointer to getMessageOpts and sets its MessageFilterMin field.
func MessageFilterMinOpt(messageFilterMin uint32) GetMessageOpt {
    return func(opts *getMessageOpts) {
        opts.MessageFilterMin = messageFilterMin
    }
}

// MessageFilterMaxOpt sets the maximum message filter value for the get message options.
//
// Parameters:
//   - messageFilterMax: The maximum value of the message filter range.
//
// Returns:
//   - GetMessageOpt: A function that takes a pointer to getMessageOpts and sets its MessageFilterMax field.
func MessageFilterMaxOpt(messageFilterMax uint32) GetMessageOpt {
    return func(opts *getMessageOpts) {
        opts.MessageFilterMax = messageFilterMax
    }
}

// DoneChanOpt sets the done channel for the get message options.
//
// Parameters:
//   - doneChan: The channel to signal when message processing is complete.
//
// Returns:
//   - GetMessageOpt: A function that takes a pointer to getMessageOpts and sets its DoneChan field.
func DoneChanOpt(doneChan chan struct{}) GetMessageOpt {
    return func(opts *getMessageOpts) {
        opts.DoneChan = doneChan
    }
}