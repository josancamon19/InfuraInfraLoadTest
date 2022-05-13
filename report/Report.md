### Load Test Report

**Base Condition:** At X requests per minute, in a linear distribution, with a timeout of 2 seconds.

```ddosify -t http://127.0.0.1:8000/eth_blockNumber -l incremental -n X -d 60 -T 2```

X=30000, the first iteration of the API (not using caching) was able to make it stable until 22k requests per minute. After this 
the connection timeout errors vary a lot, from a few <= 10 up to a few hundreds.
```
  RESULT
  -------------------------------------
  Success Count:    29987 (99%)
  Failed Count:     13    (1%)
  
  Durations (Avg):
  Server Processing    :0.2348s
  Total                :0.2348s
  
  Status Code (Message) :Count
  200 (OK)    :29987
  
  Error Distribution (Count:Reason):
  13     :connection timeout
```
   
---


### Tests using a caching layer

**Clarification:** *This caching layer is not 100% correct, it is a simplification. The right way would be using a
subscription websocket connection to the infura API and update the data depending on those notifications.*

---

Previous test were considerably expensive to run because every API call required an Infura call, and the 100k per day limit was causing trouble.

Thus, why not cache the Infura API responses and reuse those cached values? (A new ethereum block is generated every 12-14 seconds, also, some requests like 
```eth_getTransactionByBlockNumberAndIndex``` are historical, thus they will never change)

Added a caching layer using [redis](https://redis.io/). In-code Redis utils [here](../utils/redis.go)

From now on, the API will request new Infura API data every 12 seconds, all requests after and before next 12 seconds
will only be accessing to redis cached data.

 At 100k requests per minute. A minimum number of requests failed.
```
RESULT
-------------------------------------
Success Count:    99963 (99%)
Failed Count:     37    (1%)

Durations (Avg):
  Server Processing    :0.0189s
  Total                :0.0189s

Status Code (Message) :Count
  200 (OK)    :99963

Error Distribution (Count:Reason):
  37     :connection reset by peer
```

---

### Using a caching layer before we execute the router
Fiber has a middleware option for handling [caching](https://docs.gofiber.io/api/middleware/cache).
It is pretty simple and doesn't have a lot of configurations, still is enough for the task.
```go
 app.Use(cache.New(cache.Config{Expiration: 10 * time.Second}))
```

**How's this caching layer compared to the redis one?**
- ❌This one will not work in a serverless architecture.
- ❌This one has not the level of customization the redis has.
- ✅ This one is faster, as it's basically accessing a map[] object, also happens even before the route is called.
- ❌This one is taking memory directly from the API resources.

**How a request goes through?**
1. The first time, the request will ignore both caching options, and will ping the Infura API.
2. Next time, the request will go to the fiber middleware caching, after 10s, will be on (1) again
   - So, in the case of the ```eth_blockNumber``` method, it will never use the redis stuff.
   - But in case of historical data endpoints like ```eth_getTransactionByBlockNumberAndIndex```, after 10 seconds, 
     it will ignore the fiber/middleware cache, and will go to the redis cache.
     
After this setup we are now capable of going up to 250k requests per minute with an avg duration
per request of ```0.0137s```, but from 200k those results vary a lot, again, it could go up to a few hundred errors, which 
is nothing dramatic but still not ideal. Thus, I prefer keeping the max stable rpm to < 200k.

```
RESULT
-------------------------------------
Success Count:    249880 (99%)
Failed Count:     120   (1%)

Durations (Avg):
  Server Processing    :0.0136s
  Total                :0.0137s

Status Code (Message) :Count
  200 (OK)    :249880

Error Distribution (Count:Reason):
  120     :connection reset by peer
```

---

**What comes next?**

As of right now we only have `connection reset by peer` errors which is basically that the API dropped the connection 
or this was closed. 

- Explore what's the breaking point for this error starting to appear, so far it looks that happens mainly when 
  there's a sudden increase on the number of concurrent requests.
- Plot the ```ddosify``` outputs overtime, per second, and try finding new insights.
- Explore on fiber configuration object, params like `PreFork`, `Concurrency`.
- When concurrent clients call the API, and it's not time for the cache, there are many concurrent calls to the API, 
  we could synchronize the different threads/clients, making a single call to the Infura API and reusing that one
  within the other threads.