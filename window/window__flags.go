package window

type WindowDisplayFlag int32
type ClassStyleFlag int32
type WindowManagerFlag int32

const (
	// Window Display Flags
	WindowDisplayFlagShow WindowDisplayFlag = 1 << iota // Show window
	WindowDisplayFlagHide
	WindowDisplayFlagMaximize
	WindowDisplayFlagMinimize

	// Class Style Flags
	ClassStyleFlagByteAlignClient    ClassStyleFlag = 0x1000
	ClassStyleFlagByteAlignWindow    ClassStyleFlag = 0x2000
	ClassStyleFlagClassDeviceContext ClassStyleFlag = 0x0040
	ClassStyleFlagDoubleClicks       ClassStyleFlag = 0x0008
	ClassStyleFlagDropShadow         ClassStyleFlag = 0x00020000
	ClassStyleFlagGlobalClass        ClassStyleFlag = 0x4000
	ClassStyleFlagHeightRedraw       ClassStyleFlag = 0x0002
	ClassStyleFlagDisableClose       ClassStyleFlag = 0x0200
	ClassStyleFlagAllocDeviceContext ClassStyleFlag = 0x0020
	ClassStyleFlagClipParent         ClassStyleFlag = 0x0080
	ClassStyleFlagSaveBits           ClassStyleFlag = 0x0800
	ClassStyleFlagWidthRedraw        ClassStyleFlag = 0x0001
	ClassStyleFlagClipChildren       ClassStyleFlag = 0x2000000

	// Window Manager Flags
	WindowManagerFlagDestroy WindowManagerFlag = 0x0002
)
