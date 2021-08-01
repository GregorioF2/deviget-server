# Deviget Transparent Cache

Solution for transparent cache in deviget challenge.

## Project setup

```bash
 sudo docker build . -t deviget_server
 sudo docker run --network host -it deviget_server
```


## Configuration variables
Port where server should be listen to.
```
(optional) SERVER_PORT      = 3000
(optional) CACHE_MAX_TIME   = 60 // seconds
```

## Endpoints 
* **URL**

  `/prices?codes=<code_1>,<code_2>,<code_3>,...`

* **Method:**

  `GET` 
  
*  **URL Params**
  
   **Required:**
   `codes`:  List of itemcodes to request price.

* **Success Response:**

  * **Code:** 200 
    **Content:** `[ <price_1>,<price_2>,<price_3>,... ]`

## Decisions taken

- Replaced the cache price map with the one declared in the package sync (https://pkg.go.dev/sync#Map). Since this one is calculated to be much more performant, since it is optimized for two particular uses. As it says on the documentation is optimized for *when the entry for a given key is only ever written once but read many times*. That seemed to fit in the cache schema
- Replaced the type of the values in the price map with a custom entry. Since I thought it was the best solution to keep track of each entry time
```
type cacheEntry struct {
	creationTime time.Time
	value        float64
}
```

- Added the initialization of a mock service in the controller in order to test the endpoint. Of course, if there were a real service this would not be part of the solution.
 
- Added two new tests. One to make sure that if in a multiple query, one item fails, the answer of all is an error.
And another one to make sure that no Error resolved by the service is cached, since it is a common unexpected behavior in caches.

