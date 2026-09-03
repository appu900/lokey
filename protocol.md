
string  +string\r\n
bulkstrings  starts with $ follwoed by number of bytes followed by \r\n actuall string followed by \r\n can be used to store any binary data , so we are strog n which is the length of the binary 
for example emtry string  $0\r\n\r\n
null value $-1\r\n
integers starts with ':'  :1234\r\n
errors starts with '-' followed by the message followed by \r\n  -ERR unknown command\r\n


arrays :
   array starts with a '*'
   followed by number of elements 
   follwoed by '\r\n'
   followed by encoed elements 


   supose we have to store ["tony",3000,"stark"]

   *3/r/n
   $4\r\tony\r\n
   :3000\r\n
   $5\r\nstark\r\n