# MP4 Design Docs

## Maple-Juice Framework Design



## Definition

## Code Structure

## Todo

```
1. How to run specific functions designated by client (e.g. WordCount.go)?
    `Put the exec file into SDFS`

2. Should we put sdfs_src_directory into SDFS?
    `Client contact RM, which should handle this`

3. How does RM partition the Mapper task to AM?
    `Pending: First, count the total lines of data, get the workload for each AM; Second, use buffer to cache data and send to corresponding AMs`

4. How does each ApplicationMaster (AM) get file?
    `RPC or TCP Connection to get buffer then write to a local File`

5. What's the interface of Mapper & Reducer?
    `Mapper: func Mapper(fd *FIle) pair<string, string>
     Reducer: func Mapper(keyValPair pair<string, string>) pair<string, string>`

6. After Mapper Phase, what does RM do to collect "key-value" pairs?
    `Iterate through AMs and ask them to send its intermediate files to SDFS`

7. 
```