# Test parallel requests
seq 2000 | parallel -n0 -j200 "curl 127.0.0.1:8082"
