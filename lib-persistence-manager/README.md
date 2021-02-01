# lib-persistence-manager  
  
lib-persistence-manager is a library that provides an interface for Redis communication.  
Connection interface creates a connection pool, which will be used to interact with the Redis to perform CRUD operation.  
## Indexing  
lib-persistence-manager uses the Redis secondary index for indexing the resources to support search and filter capability. Currently in odimra BMC subordinate resources, Event, and Device subscriptions are indexed.
