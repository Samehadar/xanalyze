package main

import (
	"os"

	"github.com/chrislusf/gleam/distributed"
	"github.com/chrislusf/gleam/flow"
)

func main() {
	f := flow.New()
	f.TextFile("items.txt").FlatMap(`
	    function(line)
	        return line:gmatch("%w+")
	    end
	`).Map(`
	    function(word)
	        return word, 1
	    end
	`).ReduceBy(`
	    function(x, y)
	        return x + y
	    end
	`).Fprintf(os.Stdout, "%s,%d\n")

	// distributed mode
	f.Run(distributed.Option())
	f.Run(distributed.Option().SetMaster("master_ip:45326"))

}
