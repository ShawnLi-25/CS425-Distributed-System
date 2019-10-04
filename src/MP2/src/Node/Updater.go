package node

import (
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
			} else if updateQuery.queryType == 1 {
				MembershipList = AddNewNode(updateQuery.ID, MembershipList)
				SortMembershipList(MembershipList)
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
	log.Print("=== Current List is: ===")
	log.Print(list)
	newList := append(list, newNodeID)
	log.Print("=== New List is: ===")
	log.Print(newList)
	return newList
}

func DeleteNode(nodeID string, list []string) []string {
	for i := 0; i < len(list); i++ {
		if list[i] == nodeID {

			if i != len(list)-1 {
				list = append(list[:i], list[i+1:]...)
			} else {
				list = list[:i]
			}
		}
	}
	return list
}

// type MonitorList struct {
// 	localID   string
// 	monitorID []string
// 	Content   []string
// }

// func (*MonitorList) InitList() {
// 	var monitorList []string

// }
