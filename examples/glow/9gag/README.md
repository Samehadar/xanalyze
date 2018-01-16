

## Data Source
Reddit Posts 2016-09

https://storage.googleapis.com/data_interview/reddit_posts_2016_09_week1/reddit_posts_2016_09_week_1.csv.gz


## Prerequisites / Limitation

1. CSV file has to be ___uncompressed___ and ___on your local machine___ (no HDFS at the moment)
2. Only support ___standalone mode___. No cluster mode at the moment.

## Set up (Build From Source)

___Notes___: If you don't want to build the program from source, you can directly go to `q1/bin` or `q2/bin` to run the executable file. The following instruction is to install Go and build from source.

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
7. Go to `q1` or `q2` to see the instruction



## Tools

Programming Language: Go

Library: https://github.com/chrislusf/glow

## Time Reference
The following time is for reference only. The result highly depends on your machine. For Q1, there is around 200k subreddit. So, it may take longer time because the program needs to create 200k leaderboards for each subreddit.

Machine: Macbook Pro with Touchbar (i5 2.9Ghz, 16GB, SSD)

q1: 70~85 seconds

q2: 40~50 seconds

