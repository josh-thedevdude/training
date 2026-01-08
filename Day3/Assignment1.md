<pre>
REQUEST TYPE        UNDER THE HOOD

GET                 The client made a get request which reaches the server where it goes through controllers,
                    services, repositories layers that does the business logic and data fetching from database.
                    Then it returns the response to the client with appropriate status code (200)

POST                The client makes a post request to create some data on the server by sending some body with
                    the request which goes through the same cycle on the server and creates new entry in the 
                    database and returns the created data with appropraite status code (201)
</pre>
