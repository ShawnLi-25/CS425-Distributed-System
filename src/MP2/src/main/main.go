package main

import (
	"MP2/src/node"
	"flag"
)

func main() {
	isJoinPtr := flag.Bool("join", false, "join the group")
	isLeavePtr := flag.Bool("leave", false, "voluntarily leave the group")
	showListPtr := flag.Bool("membership", false, "show the current membership list")
	showIDPtr := flag.Bool("ID", false, "show self's ID")
	// configFilePtr := flag.String("config", "", "Location of Config File")

	flag.Parse()

	switch {
	case *isJoinPtr:
		node.NodeBehavior()

	}
}
