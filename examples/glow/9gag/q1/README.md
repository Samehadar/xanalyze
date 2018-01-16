## Program Design (Question 1)

This program is to use map-reduce to aggregate the CSV source. The following is the workflow:

1. Define the concurrency level based on the number of logical core of the local machine.
2. Read a CSV file and send lines via channel (Golang's pipeline)
3. Remove all the records which are empty in `subreddit` column
4. Group by user with a key pattern: `{subreddit}:::{date}:::{author}` (e.g. `pokemongo:::2016-09-01:::kittsang`)
5. Reduce users by above key and count the number of post for each user within a day
6. Break the key from `{subreddit}:::{date}:::{author}` to `{subreddit}:::{date}` in Map function
6. Group the users by the key from the result of `step 6` 
7. Order the number of user by the number of submission with descending order
8. Output the result as CSV file


## Result
The result is already uploaded called ___result.zip___ in the project folder. You can check the result directly if you want. (~200k subreddit)


## Run The Program

There are two different ways to run the program.
1. Run the executable file (Easiest)
2. Build from source and run it


### 1. Run The Executable (Easiest)
Run the executable file inside `bin` folder and run your OS version

`input`: The absolute file path of the CSV data source (e.g. /Users/name/input.csv)

Mac
```
./q1-darwin-amd64 -input [filePath]
```

Linux
```
./q1-linux-amd64 -input [filePath]
```

Windows
```
q1-windows-amd64.exe -input [filePath]
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
7. Go to the project directory `$HOME/go/src/9gag/q1` and install the dependencies by the following command
```
go get
```

8. Build and run the source under the project folder
`input`: The absolute file path of the CSV data source (e.g. /Users/name/input.csv)
```
go run main.go -input [filePath]
```
