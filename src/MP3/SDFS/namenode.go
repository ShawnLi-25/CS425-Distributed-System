package sdfs

import(
	"fmt"
	Mem "../Membership"
)

type Namenode struct{
	Filemap map[string][]string = make(map[string][]string)
}

///////////////////////////////////RPC Methods////////////////////////////
/*
	Given a request, return response containing a list of all Datanodes who has the file
*/
func (n *Namenode) GetDatanodeList (req *FindRequest, resp *FindResponse) error {
	if val, ok := Filemap[FindRequest.Filename]; ok {
		return Filemap[FindRequest.Filename]
	} 
	return nil
}

/*
	Figure out the value of Filamap[sdfsfilename] (use Mmonitoring List AKA MemHBList)
	Insert pair (sdfsfilename, datanodeList) into Filemap
	Send datanodeList back to InsertResponse
*/
func (n *Namenode) InsertFile (req *InsertRequest, resp *InsertResponse) error {
	
	datanodeList := Mem.getListByRelateIndex([]int{-2,-1,1}, InsertRequest.LocalID)

	for i, datanodeID := range datanodeList {
		Filemap[InsertRequest.Filename] = append(Filemap[InsertRequest.Filename], datanodeID) 
	}
	Filemap[InsertRequest.Filename] = datanodeList

	return datanodeList
}


///////////////////////////////////Member Function////////////////////////////

//***Function: Simply add a new entry into Filemap, return added key and value
func (n *Namenode) Add(string nodeID, string sdfsfilename) {

	
}

func (n *Namenode) Delete() {
	//TODO
	//delete an item from filemap by key
	//return deleted key and value
}

func (n *Namenode) Find() {
	//TODO
	//find value by key
	//return value if found or nil
}

func (n *Namenode) Update() {
	//TODO
	//modify value by key
	//return modified key and value
}


//////////////////////////////////////////Functions////////////////////////////////////////////

func RunNamenodeServer(Port string) {
	//TODO
}
