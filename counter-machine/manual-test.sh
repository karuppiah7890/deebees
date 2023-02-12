# Run make; ./counter-machine

# Initial counter value should be 0, if this script is run just after starting the server
curl -H "content-type: application/json" -i http://localhost:8080/counter

# Malformed JSON request should give bad request error
curl -X PATCH -H "content-type: application/json" -i http://localhost:8080/counter -d '{"something": "}'

# Empty request body (Malformed) JSON request should give bad request error
curl -X PATCH -H "content-type: application/json" -i http://localhost:8080/counter

# Sending random key-value pairs or any dummy data in the JSON request should not give any errors (for now)
# TODO: Maybe we gotta change this behavior and not accept random stuff the user sends.
# TODO: Maybe we should also restrict the size of the content that the user sends. The user won't
# be sending a very big number which is a lot of characters - meaning lot of bits / bytes - still
# the user won't be sending KBs, MBs, GBs, TBs of data. So, restrict content length to a good and
# acceptable value. And also ignore and give errors for extra keys, especially when no incrementBy key
# is not there for patch request, hmm
curl -X PATCH -H "content-type: application/json" -i localhost:8080/counter -d '{"something": "okay"}'

# What happens for this case? Does it ignore the value? No. It gives 400 bad request due to bad data type.
# But it gives a very wierd error like -
# "json: cannot unmarshal string into Go struct field JsonRequest.incrementBy of type int"
# Gotta make the error look good!
curl -X PATCH -H "content-type: application/json" -i localhost:8080/counter -d '{"incrementBy": "okay"}'

curl -X PATCH -H "content-type: application/json" -i localhost:8080/counter -d '{"incrementBy": 100 }'

# Now the counter value should be 100 :)
curl -H "content-type: application/json" -i http://localhost:8080/counter

curl -X PATCH -H "content-type: application/json" -i localhost:8080/counter -d '{"incrementBy": 100 }'

# Now the counter value should be 200 now :)
curl -H "content-type: application/json" -i http://localhost:8080/counter
