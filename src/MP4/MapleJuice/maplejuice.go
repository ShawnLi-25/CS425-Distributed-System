package maplejuice

import (
	"fmt"
	"strconv"
	"os"
	"io/ioutil"
	"log"

	config "./Config"
	sdfs "./SDFS"
)


func RunMapper(arg []string) {
	//Check argument
	mapperArg, ok := checkMapperArg(arg)
	if !ok{
		return
	}

	mapper  := mapperArg.maple_exe
	//N       := mapperArg.num_maples
	//prefix  := mapperArg.sdfs_intermediate_filename_prefix
	src_dir := mapperArg.sdfs_src_directoy

	//Upload maple_exe to SDFS
	sdfs.PutFileOrPutDir([]string{mapper, mapper})

	//Upload all files in src_dir to SDFS
	sdfs.PutFileOrPutDir([]string{src_dir, src_dir})

	//RPC Namenode's method "RunMapper"
	namenodeAddr := sdfs.GetNamenodeAddr()
	client := sdfs.NewClient(namenodeAddr + ":" + config.NamenodePort)
	client.Dial()
	defer client.Close()

	var res int
	if err := client.rpcClient.Call("Namenode.RunMapper", mapperArg, &res); err != nil {
	log.Println(err)
}

	return
}

func RunReducer(arg []string) {
	//Check argument
	reducerArg, ok := checkReducerArg(arg)
	if !ok{
		return
	}

	reducer      := reducerArg.juice_exe
	//N            := reducerArg.num_juices
	//prefix       := reducerArg.sdfs_intermediate_filename_prefix
	//destfilename := reducerArg.sdfs_dest_filename
	//delete_input := reducerArg.delete_input

	//Upload reducer_exe to SDFS
	sdfs.PutFileOrPurDir([]string{reducer, reducer})

	//RPC Namenode's method "RunReducer"
	namenodeAddr := sdfs.GetNamenodeAddr()
	client := sdfs.NewClient(namenodeAddr + ":" + config.NamenodePort)
	client.Call()
	defer client.Close()

	var res int
	if err := client.rpcClient.Call("Namenode.RunReducer", reducerArg, &res); err != nil {
	log.Println(err)
}

	return
}

/////////////////////////////Helper functions/////////////////////////////////

func checkMapperArg(arg []string) (sdfs.MapperArg, bool){
	if len(arg) < 4{
		fmt.Println("Usage: maple <maple_exe> <num_maples> <sdfs_intermediate_filename_prefix> <sdfs_src_directory>")
		return sdfs.MapperArg{}, false
	}

	//Check if maple_exe exists
	mapper  := arg[0]
	if _, err := os.Stat(mapper); os.IsNotExist(err) {
		fmt.Printf("====Error: %s not found", mapper)
		return sdfs.MapperArg{}, false
	}

	//Check if N is valid
	N, _    := strconv.Atoi(arg[1])
	if N < 0 {
		fmt.Println("====Error: non-positive num_maples")
		return sdfs.MapperArg{}, false
	}

	prefix  := arg[2]

	//Check if src_dir exists and contains file
	src_dir := arg[3]
	if _, err := os.Stat(src_dir); os.IsNotExist(err) {
		fmt.Printf("====Error: %s not found", src_dir)
		return sdfs.MapperArg{}, false
	}
	files, err := ioutil.ReadDir(src_dir)
	if err != nil {
		log.Fatal(err)
	}
	if len(files) == 0 {
		fmt.Printf("====Error: %s doesn't contains files", src_dir)
		return sdfs.MapperArg{}, false
	}

	//Return
	return sdfs.MapperArg{mapper, N, prefix, src_dir}, true
}

func checkReducerArg(arg []string) (sdfs.ReducerArg, bool){
	if len(arg) < 5{
		fmt.Println("Usage: juice <juice_exe> <num_juices> <sdfs_intermediate_filename_prefiix> <sdfs_dest_filename> delete_input={0,1}")
		return sdfs.ReducerArg{}, false
	}

	//Check if juice_exe exists
	reducer  := arg[0]
	if _, err := os.Stat(reducer); os.IsNotExist(err) {
		fmt.Printf("====Error: %s not found", reducer)
		return sdfs.ReducerArg{}, false
	}

	//Check if N is valid
	N, _    := strconv.Atoi(arg[1])
	if N < 0 {
		fmt.Println("====Error: non-positive num_juices")
		return sdfs.ReducerArg{}, false
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

	return sdfs.ReducerArg{reducer, N, prefix, destfilename, delete_input}, true
}
