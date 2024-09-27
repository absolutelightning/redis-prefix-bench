## Benchmarks Redis keys prefix search performances

### Using Scan Command

#### Setup
10 Million Keys
```python
import redis
r = redis.Redis(host='localhost', port=6379)
p = r.pipeline()
kv = {}
# Add 10 million keys to the sorted set
for i in range(1, 10000001):
    key = f"user:{i}"
    p.set(key, i)
p.execute()

print("Insertion complete.")
```

### Using Redis Module FT.SEARCH
#### Setup

Redis search docker image runs on port 6377

```python
import redis

# Connect to Redis on port 6377
r = redis.Redis(host='localhost', port=6377)
p = r.pipeline()

# Insert 10 million keys
for i in range(1, 10000001):
    p.hset(f"user:{i}", "name", f"User{i}")

p.execute()

print("Insertion complete.")
```

Run `redis-cli -p 6377`
```text
FT.CREATE idx:user ON hash PREFIX 1 "user:" SCHEMA name TEXT SORTABLE
FT.CONFIG SET MAXEXPANSIONS 1000
FT.CONFIG SET TIMEOUT 10000000
```

