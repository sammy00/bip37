package bloom

import "errors"

// ErrUninitialised signals the filter isn't initialized properly
var ErrUninitialised = errors.New("filter isn't initialised yet")
