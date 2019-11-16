package maplejuice

import {
	"fmt"
	"os"
}

type MapReduce interface{
	Map(string) map[string]string
	Reduce(map[string]string) map[string]string 
}
