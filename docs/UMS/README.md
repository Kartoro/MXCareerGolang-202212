User Management System
===========================

Environments
------------

* Golang(v1.14+)
* Mysql(v5.7+)
* Redis(v4.0+)

Flow Chart
------------
User Management System
===========================

Environments
------------

* Golang(v1.14+)
* Mysql(v5.7+)
* Redis(v4.0+)

Flow Chart
------------
![](https://github.com/Kartoro/MXCareerGolang-202212/blob/main/docs/UMS/Flow%20Chart.png)

Architecture
------------
- admin // tcp server
  - dao // mysql
  - redis // cache
  - service // tcp methods including main method
- domain // common methods and structs for admin & http 
  - model // common structs for rpc proto & sql data
  - pool // connection pool
  - rpc // rpc implementation & rpc methods
  - proto // generated proto files
  - util // common tools
- http // http server with http apis
- config
- static // static files, templates & saved avatars

HTTP APIs
------------

- GET /index - homepage
- GET /profile - get user’s profile
- POST /login - login for registered user
- POST /signup - sign up a new user
- POST /nickname - update user nickname
- POST /avatar - upload or update user avatar

RPC Methods
------------

- UserServices.SignUp
- UserServices.Login
- UserServices.GetProfile
- UserServices.UpdateNickName
- UserServices.UpdateAvatar

Architecture
------------
- admin // tcp server
  - dao // mysql
  - redis // cache
  - service // tcp methods including main method
- domain // common methods and structs for admin & http 
  - model // common structs for rpc proto & sql data
  - pool // connection pool
  - rpc // rpc implementation & rpc methods
  - proto // generated proto files
  - util // common tools
- http // http server with http apis
- config
- static // static files, templates & saved avatars

HTTP APIs
------------

- GET /index - homepage
- GET /profile - get user’s profile
- POST /login - login for registered user
- POST /signup - sign up a new user
- POST /nickname - update user nickname
- POST /avatar - upload or update user avatar

RPC Methods
------------

- UserServices.SignUp
- UserServices.Login
- UserServices.GetProfile
- UserServices.UpdateNickName
- UserServices.UpdateAvatar
