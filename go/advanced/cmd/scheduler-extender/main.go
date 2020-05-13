package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	extenderv1 "k8s.io/kube-scheduler/extender/v1"
)

func main() {
	http.HandleFunc("/filter", Filter)
	http.HandleFunc("/prioritize", Prioritize)
	log.Println("Started extender")
	if err := http.ListenAndServe(":9000", http.DefaultServeMux); err != nil {
		log.Fatalln("Server closed:", err)
	}
}

func Filter(w http.ResponseWriter, r *http.Request) {
	in := new(extenderv1.ExtenderArgs)
	if err := json.NewDecoder(r.Body).Decode(in); err != nil {
		w.WriteHeader(400)
		log.Println("Filter: bad input:", err)
		return
	}

	log.Println("Filter: called for pod. namespace=", in.Pod.Namespace, "name=", in.Pod.Name)
	for _, item := range in.Nodes.Items {
		fmt.Println("Node:", item.Name, item.Status.Addresses)
		fmt.Println("Images:", item.Status.Images)
	}
	out := &extenderv1.ExtenderFilterResult{
		Nodes:       in.Nodes,
		NodeNames:   in.NodeNames,
		FailedNodes: nil,
		Error:       "",
	}
	if err := json.NewEncoder(w).Encode(out); err != nil {
		log.Println("Filter: bad output:", err)
		return
	}
}

func Prioritize(w http.ResponseWriter, r *http.Request) {
	in := new(extenderv1.ExtenderArgs)
	if err := json.NewDecoder(r.Body).Decode(in); err != nil {
		w.WriteHeader(400)
		log.Println("Filter: bad input:", err)
		return
	}

	log.Println("Filter: called for pod ", in.Pod.Namespace, "/", in.Pod.Name)

	out := make(extenderv1.HostPriorityList, len(in.Nodes.Items))
	for i, item := range in.Nodes.Items {
		out[i].Host = item.Name
		out[i].Score = 100
	}
	if err := json.NewEncoder(w).Encode(out); err != nil {
		log.Println("Filter: bad output:", err)
		return
	}
}
