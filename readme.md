# Task Description

Using only the standard library, create a Go HTTP server that on each request responds with a counter of the total number of requests that it has received during the previous 60 seconds (moving window). The server should continue to the return the correct numbers after restarting it, by persisting data to a file.

# Todo list
[X] Read from file
[X] Write to file
-> [ ] I only need the datetime stamps in the file.
[X] Answer GET request
-> [X] Elaborate the structures
[ ] Date and time
-> [ ] Query system time
-> [ ] Write as string
-> [ ] Parse string
[ ] Put it all together
-> [X] Learn how to make functions (^.^)
-> [ ] Overwrite the file incrementally or some clever diff logic.
-> [ ] Alternatively always write to the top of the file or read from the bottom.

