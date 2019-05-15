package safe

import "runtime/debug"

// TryCatch is a helper for calling method safely.
func TryCatch(cb func(), fail func()) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			if fail != nil {
				fail()
			}
		}
	}()

	cb()
}

// Call is a wrapper of TryCatch, but set fail as nil.
func Call(cb func()) {
	TryCatch(cb, nil)
}
