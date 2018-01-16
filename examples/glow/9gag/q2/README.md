## Program Design (Question 2)

This program is to use map-reduce to aggregate the CSV source. The following is the workflow:

1. Define the concurrency level based on the number of logical core in the local machine.
2. Read a CSV file and send lines via channel (Golang's pipeline)
4. Map the key: `author` and value : `created_utc`
4. Group by the key - `author`
5. Sort the time for each key group with ascending order
6. Loop though all this time inside each key group and calculate the number of consecutive days
6. Remove all users who has the number of consecutive days = 0
7. Order the number of consecutive days (descending order)
8. Output the result `{user},{the number of consecutive days}` (No header row in the result)


## Result
The result is already uploaded called ___result.output___ in project root. You can check the result directly if you want.

## Run The Program

There are two different ways to run the program.
1. Run the executable file (Easiest)
2. Build from source and run it


### 1. Run The Executable (Easiest)
Run the executable file inside `bin` folder and run your OS version

`input`: The absolute file path of the CSV data source (e.g. /Users/name/input.csv)

Mac
```
./q2-darwin-amd64 -input [filePath]
```

Linux
```
./q2-linux-amd64 -input [filePath]
```

Windows
```
q2-windows-amd64.exe -input [filePath]
```

### 2. Build from source and run it

If you want to build from source, please follow the instruction below:

1. Install Go 1.8 or above via [https://golang.org/dl/](https://golang.org/dl/). Recommend to download `.pkg` or `.msi` for installation. Otherwise, more things need to be set up.
2. Check the installation correctness by the following command:
```
go version

// go version go1.8.x darwin/amd64

```
3. Create a directory called - `go` under your home directory `$HOME/go`
4. Create a sub-directory called - `src` (i.e. `$HOME/go/src`)
5. Export `GOPATH` by the following command
```
export GOPATH=$HOME/go
```
6. Checkout the git reporsity under `$HOME/go/src`
```
git clone https://github.com/chiukit/9gag.git
```
7. Go to the project directory `$HOME/go/src/9gag/q2` and install the dependencies by the following command
```
go get
```

8. Build and run the source under the project folder
`input`: The absolute file path of the CSV data source (e.g. /Users/name/input.csv)
```
go run main.go -input [filePath]
```

