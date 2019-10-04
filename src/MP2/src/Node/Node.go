package node

//"MP2/src/node"
//"log"
//"net"

var node Node
var UpQryChan = make(chan UpdateQuery)
var MemListChan = make(chan []string)

type Node struct {
	MemList  []string
	InGroup  bool
	Sender   Sender
	Listener Listener
	Updater  Updater
}

func CreateNewNode() Node {
	var newMemList []string
	newSender := NewSender()
	newListener := NewListener()
	newIntroducer := NewIntroducer()
	newUpdater := NewUpdater()
	newNode := Node{
		MemList:  newMemList,
		Sender:   newSender,
		Listener: newListener,
		Updater:  newUpdater,
		InGroup:  false,
	}
	return newNode
}

//Called from main.go when the command is "JOIN\n"
//Create new node and run the node until LEAVE or crash
func RunNode(isIntroducer bool) {

	node = CreateNewNode()
	if !isIntroducer {
		//false for non-intro, true for intro
		go node.NodeListen(false)
	} else {
		go node.NodeListen(true)
	}
	go node.NodeSend()
	go node.NodeUpdate()
}

//Called from main.go when the command is "LEAVE\n"
//Delete the Node
func StopNode() {
}
