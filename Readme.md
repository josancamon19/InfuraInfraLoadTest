## Infura Infra Test

### Run me

1. Run ```go get```
2. Enable ENV vars for ```API_BASE_URL```, ```INFURA_PROJECT_ID```, and ```PORT```.
3. Run ```docker run --name infura-test-redis -p 6379:6379 --rm redis``` for launching a redis db instance.
3. Run the app ```go run main.go```

**Using the Docker image:**

1. Run ```docker compose up```

### Libraries/Tools/Procedure

1. Used fiber web framework for building the Rest API, fiber is for Go what express is for node.
   WHy I used it?
   - It's Built on top of go fasthttp (an alternative to Go default net/http module) which is basically a better implementation to the way net/http launches a thread per request.
   - Fasthttp is built for handling thousands of rps, low or mid-size.
   - Fiber routing is pretty easy to set up, also has a lot of configs and middlewares available.
   - I've used it for a while.
2. Used a single endpoint `InfuraHttpRequest`. To handle a good amount of methods from the infura API.
    - First pass ```/:method``` in path. e.g. ```/eth_blockNumber``` or e.g. ```/eth_getTransactionByBlockNumberAndIndex```.
    - If the method requires parameters, you can pass them (most of the time) with a query
      param ```?params=``` with the list of params separated by comma. e.g., http://127.0.0.1:8000/eth_getTransactionByBlockNumberAndIndex?params=0x5BAD55,0x1
3. The Dockerfile is a pretty simple Multi stage Docker, compiling the Go code in container A, moving the generated
   executable to container B, running the container B.
4. Added 2 caching layers, the 1st made with redis, and the 2nd as a middleware from the fiber framework. 
   (```docker-compose.yaml``` with redis db added)
5. For the load tests, the project uses [ddosify](https://github.com/ddosify/ddosify) a pretty simple yet powerful cli 
   testing tool. DETAILS ABOUT THE LOAD TESTS RESULTS [HERE](./report/Report.md)
   


   


