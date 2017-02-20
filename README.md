# Livepeer Network Visualization

Right now the Livepeer nodes report their status to a central visualization server by default. We'll remove this after the initial build, but for now, during debugging it is useful to see which nodes connect to one another and what roles they are playing for a given stream.

## Start the Visualizaiton Server

From the `streamingviz` root directory, run

    go run ./server/server.go

Start the swarm nodes with the `--viz` flag to allow them to send their info the viz server. The `--vizhost <host>` flag will specify where the server is running (defaults to `http://localhost:8585`).

    swarm --bzzaccount 3e6a791e7b0f32fcafaa5e8fe9840dec66d5daaa --datadir ./4/ --ethapi ./1/geth.ipc --verbosity 4 --maxpeers 1 --bzznetworkid 326 --port 30403 --bzzport 8503 --rtmp 1937 --viz

## Access the visualization

If the visualization server is running, you can access the visualizaiton at [http://localhost:8585?streamid=\<streamid\>] for any given stream id. Accessing it without the argument will show the entire network, but not any stream data about who is broadcasting or consuming.

Nodes will report their peer list to the server every 20 seconds.

## TODO

1. Make sure that the http requests aren't blocking
2. Add `LogDone()` events when nodes are done broadcasting, consuming, or relaying.
3. Account for peers dropping off the network. Maybe rebuild the links at certain intervals based on latest peer data? Maybe use a timeout if we haven't seen a peer in awhile.
4. Auto refresh the visualization
5. Add a dropdown of known streams to the visualization so we can inspect them all without copying and pasting the streamID.

