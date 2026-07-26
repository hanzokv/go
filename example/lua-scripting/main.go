package main

import (
	"context"
	"fmt"

	"github.com/hanzokv/go/v9"
)

func main() {
	ctx := context.Background()

	rdb := kv.NewClient(&kv.Options{
		Addr: ":6379",
	})
	_ = rdb.FlushDB(ctx).Err()

	fmt.Printf("# INCR BY\n")
	for _, change := range []int{+1, +5, 0} {
		num, err := incrBy.Run(ctx, rdb, []string{"my_counter"}, change).Int()
		if err != nil {
			panic(err)
		}
		fmt.Printf("incr by %d: %d\n", change, num)
	}

	fmt.Printf("\n# SUM\n")
	sum, err := sum.Run(ctx, rdb, []string{"my_sum"}, 1, 2, 3).Int()
	if err != nil {
		panic(err)
	}
	fmt.Printf("sum is: %d\n", sum)
}

var incrBy = kv.NewScript(`
local key = KEYS[1]
local change = ARGV[1]

local value = kv.call("GET", key)
if not value then
  value = 0
end

value = value + change
kv.call("SET", key, value)

return value
`)

var sum = kv.NewScript(`
local key = KEYS[1]

local sum = kv.call("GET", key)
if not sum then
  sum = 0
end

local num_arg = #ARGV
for i = 1, num_arg do
  sum = sum + ARGV[i]
end

kv.call("SET", key, sum)

return sum
`)
