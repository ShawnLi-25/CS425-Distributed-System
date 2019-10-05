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

		select {
		case updateQuery := <-UpQryChan:
			if updateQuery.queryType == 0 {
				MemListChan <- MembershipList
				for _, str := range MembershipList {
					fmt.Println("Updater: Current Membership List is: " + str)
				}
			} else if updateQuery.queryType == 1 {
				MembershipList = AddNewNode(updateQuery.ID, MembershipList)
				MemListChan <- MembershipList
			} else if updateQuery.queryType == 2 {
				MembershipList = DeleteNode(updateQuery.ID, MembershipList)
				MemListChan <- MembershipList
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
		SortMembershipList(MembershipList)
		log.Print("Updater: New List is: ")
		log.Print(newList)
		return newList
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
			list = append(list[:idx], list[idx+1:]...)
		} else {
			list = list[:idx]
		}
	}
	return list
}

func FindNode(list []string, nodeID string) int {
	for i := 0; i < len(list); i++ {
		if list[i] == nodeID {
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
