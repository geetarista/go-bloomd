package bloomd

type BloomdError struct {
	ErrorString string
}

func (err *BloomdError) Error() string {
	return err.ErrorString
}

var (
	errInvalidCapacity = &BloomdError{"Must provide size with probability!"}
)

func errInvalidResponse(resp string) error {
	return &BloomdError{ErrorString: "Got response: " + resp}
}

func errCommandFailed(cmd, attempt string) error {
	return &BloomdError{
		ErrorString: "Failed to send command to bloomd server: " + cmd + ". Attempt: " + attempt,
	}
}

func errSendFailed(cmd, attempts string) error {
	return &BloomdError{
		ErrorString: "Failed to send command '" + cmd + "' after " + attempts + " attempts!",
	}
}
