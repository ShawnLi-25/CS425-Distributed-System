# MP1

CS425/ECE428 MP1 (xiangl14, taipeng2)

This is CS425 MP1-Distributed Log Querier
Group Members: Li,Xiang(xiangl14)  Liu,Taipeng(taipeng2)
Please clone the repository using `git clone https://gitlab.engr.illinois.edu/xiangl14/mp1.git`

To run the server, you need to `cd server` and use the command `go run server.go` to compile "server.go". It listens on port 8888. When recieving the input from client, server will read option and pattern from the input, execute `grep` command and output the result to client automatically.

To run the client, go to directory `cd client` and use the command `go run client.go`. Then you have to give an input. The input must follow the format `OPTION PATTERN`, for instance, `-E ^[0-9]*[a-z]{5}`. The server will run `grep` command according to the input and the client will recieve the result from server, save the result into file "MP1.log" and print out the number of lines recieved from each server.

To run the unit test, you simply need to `cd client` and then `go test -v`.Then you'll see the test result with error message if it has any.
