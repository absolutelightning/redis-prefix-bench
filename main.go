package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"strings"
	"time"

	redisearch "github.com/RediSearch/redisearch-go/redisearch"
	"golang.org/x/net/context"
)

func zrangescore(min, max int) *time.Duration {
	// Connect to Redis on port 6379
	rdb := redis.NewClient(&redis.Options{
		Addr:        "localhost:6379",
		ReadTimeout: 10 * time.Minute,
	})
	ctx := context.Background()

	_, err := rdb.Get(ctx, "a").Result()

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Error creating Redis connection:", err)
		return nil
	}

	minV := strconv.Itoa(min)
	maxV := strconv.Itoa(max)

	startTime := time.Now()

	result, err := rdb.ZRangeByScore(ctx, "ssd", &redis.ZRangeBy{
		minV, maxV,
		0, 100000000000,
	}).Result()

	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	fmt.Println(strings.Join(result, "\n"))

	endTime := time.Since(startTime)

	return &endTime
}

func scan(i int) *time.Duration {
	// Connect to Redis on port 6379
	rdb := redis.NewClient(&redis.Options{
		Addr:        "localhost:6379",
		ReadTimeout: 10 * time.Minute,
	})
	ctx := context.Background()

	_, err := rdb.Get(ctx, "a").Result()

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Error creating Redis connection:", err)
		return nil
	}

	startTime := time.Now()

	command := fmt.Sprintf("user:%d*", i)
	if i == 0 {
		command = fmt.Sprintf("user:*")
	}

	result, _, err := rdb.Scan(ctx, 0, command, 1000000000).Result()

	fmt.Println(strings.Join(result, "\n"))

	endTime := time.Since(startTime)

	return &endTime
}

func ftsearch(i int) *time.Duration {
	client := redisearch.NewClient("localhost:6377", "idx:user")
	start := time.Now()
	q := redisearch.NewQuery(fmt.Sprintf("User%d*", i))
	if i == 0 {
		q = redisearch.NewQuery("User*")
	}
	q = q.SetSortBy("name", true)
	q = q.Limit(0, 1000000)
	res, _, _ := client.Search(q)
	for _, doc := range res {
		fmt.Printf("Document ID: %s, Name: %s\n", doc.Id, doc.Properties["name"])
	}
	end := time.Since(start)
	fmt.Printf("Total Results: %d\n", end)
	return &end
}

func main() {
	//Scan
	timeDuration := scan(0)
	fmt.Println(fmt.Sprintf("user:* -> %v", timeDuration))
	for i := 1; i <= 10000000; i *= 10 {
		timeDuration = scan(i)
		fmt.Println(fmt.Sprintf("user:%d* -> %v", i, timeDuration))
	}
	for i := 1; i <= 10000000; i *= 10 {
		timeDuration = zrangescore(0, i)
		fmt.Println(fmt.Sprintf("0-%d", i), timeDuration)
	}
	timeDuration = ftsearch(0)
	fmt.Println(fmt.Sprintf("user:* -> %v", timeDuration))
	for i := 1; i <= 10000000; i *= 10 {
		timeDuration = ftsearch(i)
		fmt.Println(fmt.Sprintf("user:%d* -> %v", i, timeDuration))
	}
}
