package node

import (
	"fmt"
	"log"
	"sort"
)

var MembershipList []string

type Updater struct{}

type UpdateQuery struct {
	queryType int //0-GET, 1-ADD, 2-DELETE
	ID        string
}

//Open a go routine for this function, whenever needs update, build a channel; output will be
func (u *Updater) UpdateMembershipList() {

	for {
		// ok := <- KillRoutine
		// if ok == 1 {
		// 	fmt.Println("Updater: Go Routine Closed")
		// 	return
		// }
		select {
		case updateQuery := <-UpQryChan:
			if updateQuery.queryType == 0 {
				MemListChan <- MembershipList
				fmt.Println("Updater: Current Membership Length is: " + string(MembershipList))
				for _, str := range MembershipList {
					fmt.Printf("Updater: Membership List has member: %s...\n" + str)
				}
			} else if updateQuery.queryType == 1 {
				newMemList := AddNewNode(updateQuery.ID, MembershipList)
				MemListChan <- newMemList
			} else if updateQuery.queryType == 2 {
				newMemList := DeleteNode(updateQuery.ID, MembershipList)
				MemListChan <- newMemList
			}
		}
	}
}

func SortMembershipList(list []string) []string {
	sort.Strings(list)
	return list
}

func AddNewNode(newNodeID string, list []string) []string {
	log.Println("Updater: Current List is: ")
	log.Print(list)
	if FindNode(list, newNodeID) < 0 {
		newList := append(list, newNodeID)
		SortMembershipList(newList)
		log.Print("Updater: New List is: ")
		log.Print(newList)
		MembershipList = newList
		return MembershipList
	} else {
		return []string{}
	}

}

func DeleteNode(nodeID string, list []string) []string {
	log.Println("Updater: Current List is: ")
	log.Print(list)
	var idx = FindNode(list, nodeID)
	if idx >= 0 {
		if idx != len(list)-1 {
			MembershipList = append(list[:idx], list[idx+1:]...)
		} else {
			MembershipList = list[:idx]
		}
		return MembershipList
	} else {
		return []string{}
	}
}

func FindNode(list []string, nodeID string) int {
	for i := 0; i < len(list); i++ {
		if list[i] == nodeID {
			fmt.Println("Updater: Find Node!!!!!!" + list[i])
			return i // return index
		}
	}
	return -1
}

// type MonitorList struct {
// 	localID   string
// 	monitorID []string
// 	Content   []string
// }

// func (*MonitorList) InitList() {
// 	var monitorList []string

// }
