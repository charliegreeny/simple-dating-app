# simple-dating-app
Simple Dating API

### Running Locally 

- On initial start up, **must** run `make docker-network`
  - This creates the docker-network, `simple-dating-network`, that is needed for the differing containers to communicate

- Run `make build-and-up`
  - this build the app docker image (see `./Dockerfile`)
  - And spins that image up in a container named `dating-app`, alongside a mySql image (named `dating-db`)
  - db is seeded with 5 users (see `./scripts/init_db.sql`)

- Application will be running on `localhost:8080`

**NB** Potential "Gotcha": app fails to build when both db and spin up at same time:

- The app image is depended on the mySql image starting (see `./docker-compose.yml`)
- _However_, when the mySql image starts it is not initially open to connections and the `dating-app` container will stop
- **Therefore, the fix is:**
  - wait ~1 min OR until `dating-db` prints the following log message
    
  ```
  "X Plugin ready for connections. Bind-address: '::' port: 33060, socket: /var/run/mysqld/mysqlx.sock"
  ```
    
  -  Then run `make up`

### Custom Endpoint login 

- `POST /swipe`, `POST /user/preference` & `GET /discover` all **require** `Authorization` header set with the token provided in the `/login` response
- Similarly, for the above endpoints you can add an optional `x-user-locatoin` to update the authenticated user's location: 
  - eg 
  ```json
    {
      "lat": 53.4808, // float
      "long":-2.2426  // float
  }

  ```
  
- `PATCH /user/perference`
  -  request/response body: 
  ```json 
   {"gender": "MALE",
        "minAge": 18,
        "maxAge": 27,
        "maxDistance": 50}
   ```


### Application logic 

- On start up adds all users from the db add to a cache
- On `POST` requests, e.g. `/user/create/` the app will insert into the db and asynchronously add to said cache
  - the same is true for any data updates in the app and `PATCH` request (e.g `user/preferences`)
- Equally, any "gets" throughout the app would check the cache first then the database
  - see implementations of the `app.IDGetter` interface

- Login/Auth  
  - `POST /login` creates a jwt and stores that in a cache against the user id provided in the request
  - This cache is used to get the user from the userCache and add that to `context` on requests that need to be logged in for 
    - i.e. endpoints that require an `Authorization` header (the `token` in a successful login request)
    - e.g. `GET /discovery` & `POST /swipe`

- Created a preferences object to allow for filtering 
  - Filtering done in order of: preferred gender, with min/max age, then is within max distance (km)
  - "Default" preference are added on user creation, where:
    - preferred gender: the opposite of user provided,
    - min age: 18
    - max age: null (no upper bound essentially)
    - maxDistance: 100 (km)

[//]: # (  - update preferences with `PATCH /user/preferences` with the login token provided in `Authorization` header)

### Decisions Made

- User id as an uuid string 
  - having id as an int has previously allowed for a security flaw as can be predictable (e.g. [here](https://www.youtube.com/watch?v=CgJudU_jlZ8&ab_channel=TomScott))
  - Therefore, I see it as best practice to use a none predictable string, especially an uuid, as best practice
  - Plus, a random int would open to many of the same issues of an uuid (probably more so as shorter and only 0-9 user vs uuid's alphanumeric)

- Not returning password 
  - I was already uneasy in accepting a password string in 2 requests (`POST /user/create` & `POST /login`), returning it would also compound that feeling of unease in not following best practises

- SQL database 
  - Owing to the perceived rigidity in the data model, I thought it would be best to have a SQL database. In order to allow for a clear data contract that could be unnecessarily flexible in a schemaless database

- In memory caches 
  - Used to reduce complexity and time to make app

- Updating location 
  - I thought it would be useful to allow location to be updated during any requests where the user is logged in (e.g. GET /discovery)
  - This is on the assumption there would be a mobile app using this application, therefore location can change regularly. 
  - To allow this having "middleware" that took a json body in a `x-user-location` would allow for a customer current location rather than a stored value 
  - The application would resolve to the stored value (cache or db) if the header was not provided

- Unit test
  - Cut short for time
  - I applied them where I found the most utility for them; i.e. :
    - Where it was not running as I expected and therefore added tests to make sure it robust 
    - And not using them where functionality was similar (e.g. get a user from the db vs get a preferences record fromm the db) 

### Improvements/Wants
Owing to time and keeping the app relatively simple here are a few things I would have liked to have done that would have improved the app

- Fetching all users on start up
  - Would the **1st improvement** for this app if it were to be scaled
  - One solution I can think is adding users to the cache as they log in and having a long TTL (e.g. number of days/weeks)
  - Or having an external cache (see below) that would be topped up by a worker, either attached to the database or following its own unique logic

- TTL for cached items 
  - Again cut for time. I almost went on forever and forever adding things so had to be selectively on what to implement
  - I would have used something like [this](https://github.com/ellogroup/ello-golang-cache) which was lifted from a private app at ElloGroup (prev role) I created.
  - This implements fetcher which will get the item from its source if it had expired or not in cache

- Using External caches
  - There are many uses for either using a cache, like Redis, for the caches in this app - especially the `userCache` and `jwtCache`
  - This would allow for the data to be 
    - i) accessible by other services if required (especially true of the user cache); and
    - ii) persist beyond the up-time of this service

- Non-SQL store for matches
  - Simply to the above this would allow other services to access it and persist longer than the service
  - Plus, if, for example, DynamoDB was used it would be possible to leverage Event-Driven architecture to notify the other user using AWS DynamoStreams & Lambdas.
  - Furthermore, this DynamoDB would offer the flexible to allow for other user-related events that are implemented later. So would be "user-event" data-store than simply a "matches" store

- Testing
  - End-to-end tests using, for example, cucumber and [godog](https://github.com/cucumber/godog) would offer great value in my opinion
  - It would allow to cover who user flows, for example, creating a user -> logging in -> updating preferences/location -> getting a match