package client

import (
	"testing"
)

// Need to be running a server on port 8585 to run the below.
// For now treat it as an example usage
func TestVizClient(t *testing.T) {
	host := "http://localhost:8585"
	client := NewClient("A", true, host)
	client.LogPeers([]string{"B", "C"})
	client.LogBroadcast("stream1")

	client2 := NewClient("B", true, host)
	client2.LogPeers([]string{"A", "D"})

	client3 := NewClient("D", true, host)
	client3.LogPeers([]string{"B"})
	client3.LogConsume("stream1")

	client2.LogRelay("stream1")
}
