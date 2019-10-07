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
var MonitorList []string

type Updater struct{}

type UpdateQuery struct {
	queryType int //0-GET, 1-ADD, 2-DELETE
	ID        string
}

//func UpdateMemshipList(

//Open a go routine for this function, whenever needs update, build a channel; output will be
// func (u *Updater) UpdateMembershipList() {
// 	for {
// 		select {
// 		case updateQuery := <-UpQryChan:
// 			// if updateQuery.queryType == 0 {
// 			// 	MemListChan <- MembershipList
// 			// 	log.Printf("Updater: Current Membership Length is: %d...\n ", len(MembershipList))
// 			// 	for _, str := range MembershipList {
// 			// 		log.Printf("Updater: Membership List has member: %s...\n", str)
// 			// 	}
// 			// } else
// 			if updateQuery.queryType == 1 {
// 				newMemList := AddNewNode(updateQuery.ID, MembershipList)
// 				MemListChan <- newMemList
// 				UpdateMemHBMap()

// 				for _, str := range MonitorList {
// 					log.Printf("Updater: MonitorList has member:%s...\n", str)
// 					fmt.Printf("Updater: MonitorList has member:%s...\n", str)
// 				}
// 				for NodeID, _ := range MemHBMap {
// 					log.Printf("Updater: MemHBMap has member:%s...\n", NodeID)
// 					fmt.Printf("Updater: MemHBMap has member:%s...\n", NodeID)
// 				}
// 			} else if updateQuery.queryType == 2 {
// 				newMemList := DeleteNode(updateQuery.ID, MembershipList)
// 				MemListChan <- newMemList
// 				UpdateMemHBMap()
// 				for _, str := range MonitorList {
// 					log.Printf("Updater: MonitorList has member:%s...\n", str)
// 					fmt.Printf("Updater: MonitorList has member:%s...\n", str)
// 				}
// 				for NodeID, _ := range MemHBMap {
// 					log.Printf("Updater: MemHBMap has member:%s...\n", NodeID)
// 					fmt.Printf("Updater: MemHBMap has member:%s...\n", NodeID)
// 				}
// 			}
// 		case <-KillUpdater:
// 			// ln.Close()
// 			fmt.Println("===Updater: UpdateMembershipList Leave!!")
// 			// KillRoutine <- struct{}{}
// 			return
// 		}
// 	}
// }

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

func AddNewNode(newNodeID string) []string {
	log.Println("Updater: Before ADD Current List is: ")
	log.Print(MembershipList, "\n")
	if FindNode(MembershipList, newNodeID) < 0 {
		newList := append(MembershipList, newNodeID)
		SortMembershipList(newList)
		log.Print("Updater: New List is: ")
		log.Print(newList, "\n")
		MembershipList = newList
		return MembershipList
	} else {
		return []string{}
	}

}

func DeleteNode(nodeID string) []string {
	log.Println("Updater: Before Delete the List is: ")
	fmt.Print(MembershipList, "\n")
	var idx = FindNode(MembershipList, nodeID)
	log.Printf("The Delete Node is in the positition: %d\n", idx)
	if idx != -1 {
		if idx != len(MembershipList)-1 {
			MembershipList = append(MembershipList[:idx], MembershipList[idx+1:]...)
		} else {
			MembershipList = MembershipList[:idx]
		}
		newList := make([]string, len(MembershipList))
		copy(newList, MembershipList)
		log.Print("Updater: New List is: ")
		log.Print(newList, "\n")
		// newList := MembershipList
		return newList
	} else {
		return []string{}
	}
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
