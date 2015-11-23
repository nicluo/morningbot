package morningbot

import (
  "runtime"
)

// GoSafely is a utility wrapper to recover and log panics in goroutines.
// If we use naked goroutines, a panic in any one of them crashes
// the whole program. Using GoSafely prevents this.
func (m *MorningBot) GoSafely(fn func()) {
  go func() {
    defer func() {
      if err := recover(); err != nil {
        stack := make([]byte, 1024*8)
        stack = stack[:runtime.Stack(stack, false)]

        m.log.Printf("PANIC: %s\n%s", err, stack)
      }
    }()

    fn()
  }()
}
