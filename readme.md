# Task Description

Using only the standard library, create a Go HTTP server that on each request responds with a counter of the total number of requests that it has received during the previous 60 seconds (moving window). The server should continue to the return the correct numbers after restarting it, by persisting data to a file.

## Refinement of requirements

[X] The server must be race-condition free and able to handle multiple requests.

[X] The server should be efficient when interacting with the filesystem.

[ ] Precision should not be compromised.

[ ] There should be proper tests in place to make sure that these conditions hold.

## Architecture

The pipeline of data looks like this:
```
......... f    ..........   g   ........  h  .........
Server --writes to---> datastructure -writes to-> file
       \_concurrent_/
```
Constraints:

- f:  can only append to queue
- g:  is an endomorphism that keeps old data outside
- h1: has read only access to queue
- h2: writes efficiently to file

h = h2 . h1

Procrastinating a little bit, made a better picture :)
```
--- Server Thread ---
| 1 | 2 | 3 | 4 | 5 |  :: f
\   \   \   /   /   /
.\   \   \ /   /   /.
+>\ Datastructure°/..  :: g
|..\|||||||||||../...  :: h1
|.. Write Logic ----+  :: h2
+-on start-- File <-+
```
## Philosophy

I decided to use a FIFO queue instead of PriorityQueue because it's easier and faster but there is the possibility of there being some jitter (where timestamps that are very close to each other get swapped, i.e. put in inverse order).

This means that I prioritized throughput over precision.

However I can only use the standard library and they don't offer a good queue so the throughput is somewhat dampened by that.

The user is asking at the granularity of second so the answer is approximately correct for this granularity.


# Todo list

[ ] Write tests

-> [ ] Accuracy: time, free of raceconditions

-> [ ] Stress: generate concurrent events

[ ] Improve time layout

-> [ ] Precision: Find out how much I can have

-> [ ] Efficiency: Find out how much it costs

-> [ ] Notation: Find out how to communicate it


