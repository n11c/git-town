package helpers

// removeStringChannel provides a copy of the given list of string channels
// with the channel at the given index removed.
func removeStringChannel(channels []chan string, index int) []chan string {
	channels[index] = channels[len(channels)-1]
	return channels[:len(channels)-1]
}
