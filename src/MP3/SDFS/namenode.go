package sdfs

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"reflect"

	Config "../Config"
	Mem "../Membership"
)

var namenode = new(Namenode)
var membershipList []string

type Namenode struct {
	Filemap map[string][]string //Key: sdfsFileName  Value: Arraylist of replica node
	Nodemap map[string][]string //Key: NodeID  Value: Arraylist of sdfsFileName
}

//////////////////////////////////////////Functions////////////////////////////////////////////

func RunNamenodeServer() {
	namenodeServer := rpc.NewServer()

	err := namenodeServer.Register(namenode)
	if err != nil {
		log.Fatal("Register(namenode) error:", err)
	}

	//======For multiple servers=====
	oldMux := http.DefaultServeMux
	mux := http.NewServeMux()
	http.DefaultServeMux = mux
	//===============================

	namenodeServer.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

	//=======For multiple servers=====
	http.DefaultServeMux = oldMux
	//================================

	listener, err := net.Listen("tcp", ":"+Config.NamenodePort)
	if err != nil {
		log.Fatal("Listen error", err)
	}

	fmt.Printf("===RunNamenodeServer: Listen on port %s\n", Config.NamenodePort)
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("Serve(listener, nil) error: ", err)
	}

}

//***Todo: Check if it's correct
func updateNameNode(newMemList []string) {
	for {
		var addList, deleteList []string
		mapEq := reflect.DeepEqual(newMemList, membershipList)
		if !mapEq {
			for newIdx, oldIdx := 0, 0; newIdx < len(newMemList) && oldIdx < len(membershipList); {
				if newMemList[newIdx] == membershipList[oldIdx] {
					newIdx++
					oldIdx++
				} else {
					//
					if newMemList[newIdx] < membershipList[oldIdx] {
						addList = append(addList, newMemList[newIdx])
						log.Printf("===New Added Node:%s\n", newMemList[newIdx])
						newIdx++
					} else {
						deleteList = append(deleteList, membershipList[oldIdx])
						log.Printf("===Deleted Node:%s\n", membershipList[oldIdx])
						oldIdx++
					}
				}
			}
		}
		updateMap(addList, deleteList)
		reReplicate(deleteList)
		membershipList = newMemList
	}
}

//***Todo: Update two essential map
func updateMap(addList []string, deleteList []string) {

}

//Todo: Rereplicate files for deleting Node
func reReplicate(deleteList []string) {
	reFileSet := make(map[string]bool)
	for _, nodeID := range deleteList {
		for _, fileName := range namenode.Nodemap[nodeID] {
			if _, ok := reFileSet[fileName]; !ok {
				reFileSet[fileName] = true
				//***Todo: Replicate from sdfsfile?
				fmt.Printf("===Re-replicate file: %s!!!\n", fileName)
				PutFile([]string{fileName, fileName}, false)
			}
		}

	}
}

///////////////////////////////////RPC Methods////////////////////////////
/*
	Given a request, return response containing a list of all Datanodes who has the file
*/

/*TODO Implement GetDatanodeList
func (n *Namenode) GetDatanodeList(req FindRequest, resp *FindResponse) error {
	resp.DatanodeList = []string{"fa19-cs425-g73-01.cs.illinois.edu",
				     "fa19-cs425-g73-02.cs.illinois.edu",
				     "fa19-cs425-g73-03.cs.illinois.edu"}
	return nil
}

TODO Implement InsertFile
func (n *Namenode) InsertFile(req InsertRequest, resp *InsertResponse) error {
	resp.DatanodeList = []string{"fa19-cs425-g73-01.cs.illinois.edu",
				     "fa19-cs425-g73-02.cs.illinois.edu",
				     "fa19-cs425-g73-03.cs.illinois.edu"}
	return nil
}
*/

func (n *Namenode) GetDatanodeList(req *FindRequest, resp *FindResponse) error {

	if val, ok := n.Filemap[FindRequest.Filename]; ok {
		resp.DatanodeList = n.Filemap[FindRequest.Filename]
	} else {
		resp.DatanodeList = []string{}
	}
	return nil
}

/*
	First time for put original file (Assign to Mmonitoring List AKA MemHBList)
	Insert pair (sdfsfilename, datanodeList) into Filemap
	Send datanodeList back to InsertResponse
*/

func (n *Namenode) InsertFile(req InsertRequest, resp *InsertResponse) error {

	datanodeList := Mem.GetListByRelateIndex([]int{-2, -1, 1}, InsertRequest.LocalID)

	for i, datanodeID := range datanodeList {
		fmt.Fprintf("**namenode**: Insert sdfsfile: %s to %s from %s\n", InsertRequest.Filename, datanodeID, InsertRequest.LocalID)
		log.Printf("**namenode**: Insert sdfsfile: %s to %s from %s\n", InsertRequest.Filename, datanodeID, InsertRequest.LocalID)
		n.Filemap[InsertRequest.Filename] = append(n.Filemap[InsertRequest.Filename], datanodeID)
		n.Nodemap[datanodeID] = append(n.Nodemap[datanodeID], InsertRequest.Filename)
	}
	// n.Filemap[InsertRequest.Filename] = datanodeList

	resp.DatanodeList = datanodeList
}

//TODO
//Note: Map operation is not required to be implemented.
//If we do, please implement them into FUNCTION, NOT METHOD.
//The reason is that class Namenode is registered in RPC.
//All methods of Namenode MUST have a standard format like
//func (a Type) method([Valuable of Explicit Type], [Pointer of Explicit Type]) error{}

/*
func (n *Namenode) Add(nodeID string, sdfsfilename string) {
	return
}

func (n *Namenode) Delete() {
	//TODO
	//delete an item from filemap by key
	//return deleted key and value
	return
}

func (n *Namenode) Find() {
	//TODO
	//find value by key
	//return value if found or nil
	return
}

func (n *Namenode) Update() {
	//TODO
	//modify value by key
	//return modified key and value
	return
}


*/
