package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/livepeer/streamingviz"
)

var networkDB *streamingviz.Network

func init() {
	networkDB = streamingviz.NewNetwork()
}

func handler(w http.ResponseWriter, r *http.Request) {
	streamID := r.URL.Query().Get("streamid")

	abs, _ := filepath.Abs("./server/static/index.html")
	view, err := template.ParseFiles(abs)

	//data := getData()
	//network := getNetwork()
	//data := networkToData(network, "teststream")

	data := networkToData(networkDB, streamID)

	if err != nil {
		fmt.Fprintf(w, "error: %v", err)
	} else {
		view.Execute(w, data)
	}
}

func handleJson(w http.ResponseWriter, r *http.Request) {
	abs, _ := filepath.Abs("./server/static/data.json")
	view, _ := ioutil.ReadFile(abs)
	fmt.Fprintf(w, "%s", view)
}

func handleEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "Error. You must POST events")
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	var event map[string]interface{}
	if err := json.Unmarshal(body, &event); err != nil {
		fmt.Fprintf(w, "Error unmarshalling request: %v", err)
		return
	}

	eventName := event["name"].(string)
	node := event["node"].(string)

	switch eventName {
	case "peers":
		peers := event["peers"].([]interface{})
		peerList := make([]string, 0)
		for _, v := range peers {
			peerList = append(peerList, v.(string))
		}
		networkDB.ReceivePeersForNode(node, peerList)
	case "broadcast":
		streamID := event["streamId"].(string)
		fmt.Println("Got a BROADCAST event for", node, streamID)
		networkDB.StartBroadcasting(node, streamID)
	case "consume":
		streamID := event["streamId"].(string)
		fmt.Println("Got a CONSUME event for", node, streamID)
		networkDB.StartConsuming(node, streamID)
	case "relay":
		streamID := event["streamId"].(string)
		fmt.Println("Got a RELAY event for", node, streamID)
		networkDB.StartRelaying(node, streamID)
	case "done":
		streamID := event["streamId"].(string)
		fmt.Println("Got a DONE event for", node, streamID)
		networkDB.DoneWithStream(node, streamID)
	case "default":
		fmt.Fprintf(w, "Error, eventName %v is unknown", eventName)
		return
	}
}

func main() {
	http.HandleFunc("/data.json", handleJson)
	http.HandleFunc("/event", handleEvent)
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8585", nil)
}

func getNetwork() *streamingviz.Network {
	sID := "teststream"
	network := streamingviz.NewNetwork()

	// Set up peers
	network.ReceivePeersForNode("A", []string{"B", "D", "F"})
	network.ReceivePeersForNode("B", []string{"D", "E"})
	network.ReceivePeersForNode("C", []string{"I", "F", "G"})
	network.ReceivePeersForNode("E", []string{"I", "H"})
	network.ReceivePeersForNode("F", []string{"G"})
	network.ReceivePeersForNode("G", []string{"H"})
	network.ReceivePeersForNode("H", []string{"I"})

	network.StartBroadcasting("A", sID)
	network.StartConsuming("I", sID)
	network.StartConsuming("G", sID)
	network.StartRelaying("F", sID)
	network.StartRelaying("C", sID)
	network.DoneWithStream("B", sID)

	return network
}

func networkToData(network *streamingviz.Network, streamID string) interface{} {
	/*type Node struct {
		ID string
		Group int
	}

	type Link struct {
		Source string
		Target string
		Value int
	}*/

	res := make(map[string]interface{})
	nodes := make([]map[string]interface{}, 0)

	for _, v := range network.Nodes {
		nodes = append(nodes, map[string]interface{}{
			"id":    v.ID,
			"group": v.GroupForStream(streamID),
		})
	}

	links := make([]map[string]interface{}, 0)

	for _, v := range network.Links {
		links = append(links, map[string]interface{}{
			"source": v.Source.ID,
			"target": v.Target.ID,
			"value":  2, //v.Value[streamID],
		})
	}

	res["nodes"] = nodes
	res["links"] = links

	b, _ := json.Marshal(res)
	fmt.Println(fmt.Sprintf("The output network is: %s", b))

	var genResult interface{}

	json.Unmarshal(b, &genResult)
	return genResult
}

func getData() map[string]interface{} {
	return map[string]interface{}{
		"nodes": []map[string]interface{}{
			{
				"id":    "A",
				"group": 1,
			},
			{
				"id":    "B",
				"group": 2,
			},
		},
		"links": []map[string]interface{}{
			{
				"source": "A",
				"target": "B",
				"value":  1,
			},
		},
	}
}
