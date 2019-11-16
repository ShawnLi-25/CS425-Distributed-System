package maplejuice

import (
	"fmt"
	"strconv"
	"os"
	"ioutil"
)

type MapperArg struct{
	maple_exe string
	num_maples int
	sdfs_intermediate_filename_prefix string
	sdfs_src_directory string
}

type ReducerArg struct{
	juice_exe string
	num_juices int
	sdfs_intermediate_filename_prefix string
	sdfs_dest_filename string
	delete_input bool
}

func RunMapper(arg []string) {
	//Check argument
	mapperArg, ok := checkMapperArg(arg)
	if !ok{
		return
	}

	mapper  := mapperArg.maple_exe
	N       := mapperArg.num_maples
	prefix  := mapperArg.sdfs_intermediate_filename_prefix
	src_dir := mapperArg.sdfs_src_directory

	//Figure out buffer size
	//TODO

	//Split all date into buffers
	//Write buffer into file before calling mapper
	//TODO

	//Wait all mappers
	//TODO

	//Return
}

func RunReducer(arg []string) {
	//Check argument
	reducerArg, ok := checkReducerArg(arg)
	if !ok{
		return
	}

	reducer      := reducerArg.juice_exe
	N            := reducerArg.num_juices
	prefix       := reducerArg.sdfs_intermediate_filename_prefix
	destfilename := reducerArg.sdfs_dest_filename
	delete_input := reducerArg.delete_input

	//Figure out buffer size
	//TODO
}

/////////////////////////////Helper functions/////////////////////////////////

func checkMapperArg(arg []string) (MapperArg, bool){
	if len(arg) < 4{
		fmt.Println("Usage: maple <maple_exe> <num_maples> <sdfs_intermediate_filename_prefix> <sdfs_src_directory>")
		return MapperArg{}, false
	}

	//Check if maple_exe exists
	mapper  := arg[0]
	if _, err := os.Stat(mapper); os.IsNotExist(err) {
		fmt.Printf("====Error: %s not found", mapper)
		return MapperArg{}, false
	}

	//Check if N is valid
	N, _    := strconv.Atoi(arg[1])
	if N < 0 {
		fmt.Println("====Error: non-positive num_maples")
		return MapperArg{}, false
	}

	prefix  := arg[2]

	//Check if src_dir exists and contains file
	src_dir := arg[3]
	if _, err := os.Stat(src_dir); os.IsNotExist(err) {
		fmt.Printf("====Error: %s not found", src_dir)
		return MapperArg{}, false
	}
	files, err := ioutil.ReadDir(src_dir)
	if err != nil {
		log.Fatal(err)
	}
	if len(files) == 0 {
		fmt.Printf("====Error: %s doesn't contains files", src_dir)
		return MapperArg{}, false
	}

	//Return
	return MapperArg{mapper, N, prefix, src_dir}, true
}

func checkReducerArg(arg []string) (ReducerArg, bool){
	if len(arg) < 5{
		fmt.Println("Usage: juice <juice_exe> <num_juices> <sdfs_intermediate_filename_prefiix> <sdfs_dest_filename> delete_input={0,1}")
		return ReducerArg{}, false
	}

	//Check if juice_exe exists
	reducer  := arg[0]
	if _, err := os.Stat(reducer); os.IsNotExist(err) {
		fmt.Printf("====Error: %s not found", reducer)
		return ReducerArg{}, false
	}

	//Check if N is valid
	N, _    := strconv.Atoi(arg[1])
	if N < 0 {
		fmt.Println("====Error: non-positive num_juices")
		return ReducerArg{}, false
	}

	prefix  := arg[2]
	//TODO what if no sdfsfile has matching prefix?

	destfilename := arg[3]

	var delete_input bool
	if arg[4] == "delete_input=0" || arg[4] == "0" {
		delete_input = false
	}else if arg[4] == "delete_input=1" || arg[4] == "1" {
		delete_input = true
	}else {
		//By default
		delete_input = false
	}

	return ReducerArg{reducer, N, prefix, destfilename, delete_input}, true
}

//TODO: Put this function into Config.go
//Count all contents size of a directoru, not including sub-dir
func getDirSize(dirpath string) int64{
	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		log.Fatal(err)
	}

	
}
