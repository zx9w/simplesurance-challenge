# Task Description

Using only the standard library, create a Go HTTP server that on each request responds with a counter of the total number of requests that it has received during the previous 60 seconds (moving window). The server should continue to the return the correct numbers after restarting it, by persisting data to a file.

# Todo list
[X] Read from file
[ ] Write to file
[ ] Answer GET request
[ ] Measure time
[ ] Put it all together

I have to remember to put timestamps on increments in the file. Overwrite the stuff older than one minute when rewriting.
