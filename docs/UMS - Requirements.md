# Requirements

## Features

Implement a user manager system for user to login and edit their profiles. Users can login on web page.

After a user inputs his/her username and password, the backend system will authenticate his/her identity. If the login is successful, the relevant user information is displayed, otherwise an error message will be shown.

After a successful login, a user can edit the following information:

- Upload a picture as his/her profile picture
- Change his/her nickname (support unicode characters with utf-8 encoding)
- User information includes: username (cannot be changed), nickname, profile picture. For test, the initial user data can be directly insert into database. Make sure there are at least 10 million user accounts in the test database.

## Performance

Supports 1000 concurrent http requests 

### With cache layer

Basic requirement: 
- supports 6000 login requests per second from at least 200 unique users
- no socket/timeout error
- supports 20k login requests per second (Optional, you can try reach this target if have time) 

### Without cache layer

Basic requirement: 
- supports 3000 login requests per second from at least 200 unique users
- no socket/timeout error
- supports 10k login requests per second (Optional)

### Stress Test Environment

OS: MacOS or Linux VM or WIN10
DB: MySQL 5.7
Client: Chrome and Firefox


## Design Constraint

- Separate HTTP server and TCP server and put the main logic on TCP server
- Backend authentication logic should be done in the TCP server
- User information must be stored in a MySQL database. Connect by MySQL Go client.
- Use standard library whenever possible. Build the connection pool library yourself.
- Web interface will not directly connect to MySQL. For each HTTP request, web interface will send a TCP request to the TCP server, which will handle business logic and query the database.
- The request and response should be in protobuf format. https://developers.google.com/protocol-buffers

## Design Considerations

- Robustness
- Security
- Performance

## References
Go: http://golang.org

Coding style: https://github.com/golang/go/wiki/CodeReviewComments

Testing: https://golang.org/pkg/testing/

Profiling: http://blog.golang.org/profiling-go-programs

Go Web application example: https://golang.org/doc/articles/wiki/

Go editor/IDE

https://www.jetbrains.com/go/ (Recommended)
https://code.visualstudio.com/
https://github.com/fatih/vim-go
https://github.com/dominikh/go-mode.el
https://github.com/DisposaBoy/GoSublime
https://github.com/visualfc/liteide
You can apply for software license via IT-Centre (https://it-center.seagroup.com/create-ticket)
MySQL client: https://github.com/go-sql-driver/mysql

Redis: http://redis.io

Redis Client: https://github.com/go-redis/redis

Protobuf: https://developers.google.com/protocol-buffers

MySQL Database Design Guide
