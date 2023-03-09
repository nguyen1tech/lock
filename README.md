# Lock

## Problem:
Implement a lock service

### 1. Understand the requirements and establish the design scope
    - Support 2 operations: Acquire a lock and Release a lock
    - Should work in distributed environment
    - Performant, high availibility, consistency
    - A lock automatically releases when the expiration has reached
### 2. High level design
    - Client --> Lock Service
    
    - Acquire(key): acquires a lock with the given key 
    - Release(key): releases a lock held on the given key

    - Usecases:
    + client1 acquires a lock on "ABC" key -> client1 holds a lock on "ABC" key
    + client2 acquires a lock on "ABC" key -> fails dues to client1 is holding a lock on "ABC" key
    + client3 acquires a lock on "CDF" key -> lock acquired on "CDF" key
    + client1 finishes the work and release the lock on "ABC" key
    + client2 acquires a lock on "ABC key -> client2 holds a lock on "ABC" key
### 3. Deep dive
    - What if client1 crashes while holding a lock on "ABC" key -> ?
    - Should we queue clients who are trying to acquire a lock on a specified key as another alternative -> ?
    - If the Lock Service crashes all lock information currently stores in the memory will be lost -> Persist to disk?
### 4. Wrap up

### 5. Test
```
-- server
docker run --name redis-instance -p 6379:6379 -d redis:7.0.8
-- cli
docker run -it --rm redis redis-cli -h redis-instance
```