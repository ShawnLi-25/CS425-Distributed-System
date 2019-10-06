package node

import (
	"fmt"
	"log"
	"sort"
	"time"

	msg "../Helper"
)

var MembershipList []string
var MemHBMap map[string]time.Time = make(map[string]time.Time)

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
				log.Printf("Updater: Current Membership Length is: %d...\n ", len(MembershipList))
				for _, str := range MembershipList {
					log.Printf("Updater: Membership List has member: %s...\n", str)
				}
			} else if updateQuery.queryType == 1 {
				newMemList := AddNewNode(updateQuery.ID, MembershipList)
				MemListChan <- newMemList
				UpdateMemHBMap()
			} else if updateQuery.queryType == 2 {
				newMemList := DeleteNode(updateQuery.ID, MembershipList)
				MemListChan <- newMemList
				UpdateMemHBMap()
			}
		}
	}
}

//Use MembershipList to update the key in MemHBMap(NodeID, Time)
func UpdateMemHBMap() {
	var newMemHBMap map[string]time.Time = make(map[string]time.Time)
	MemHBList := msg.GetMonitoringList(MembershipList, LocalAddress)

	if len(MemHBMap) == 0 {
		for _, c := range MemHBList {
			MemHBMap[c] = time.Now()
		}
	} else {
		for _, c := range MemHBList {
			if LastTime, ok := MemHBMap[c]; ok {
				newMemHBMap[c] = LastTime
			} else {
				newMemHBMap[c] = time.Now()
			}
		}
		MemHBMap = newMemHBMap
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
		return newList
	} else {
		return []string{}
	}

}

func DeleteNode(nodeID string, list []string) []string {
	fmt.Println("Updater: Current List is: ")
	fmt.Print(list, "\n")
	fmt.Println("Updater: Delete Node ID is: " + nodeID)
	var idx = FindNode(list, nodeID)
	fmt.Printf("Updater: Find Index is: %d", idx)
	if idx != -1 {
		if idx != len(list)-1 {
			MembershipList = append(list[:idx], list[idx+1:]...)
		} else {
			MembershipList = list[:idx]
		}
		newList := make([]string, len(MembershipList))
		copy(newList, MembershipList)
		// newList := MembershipList
		return newList
	} else {
		fmt.Println("What the fuck!!!!!!!")
		return []string{}
	}
}

func FindNode(list []string, nodeID string) int {
	fmt.Printf("Updater: Current Length is: %d", len(list))
	for i := 0; i < len(list); i++ {
		if list[i] == nodeID {
			fmt.Printf("Updater: Find Node %s at position %d!!!!!!\n", list[i], i)
			return i // return index
		} else {
			fmt.Println("Updater: No Match Node!!!!!!" + list[i])
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
