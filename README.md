# IES PIK Church Management System

## Overview
This is an open source alternative for church management system applications.
The goal of this project is to provide cheaper alternatives to small church that can't afford to use paid software to manage the church.


## Environment variable

```
SERVICE_NAME=CMS
SERVICE_MODULES=AUTH,PEOPLE,EVENTS

EMAIL_ADDRESS=<Email Sender address for OTP. must be of google account>;
EMAIL_PASSWORD=<Email Password for OTP. must be of google account>;
ADMIN_USER=<Email of system admin>;
ADMIN_PASSWORD=<Password of system admin>;


CMS_MONGO_URL=localhost:27017
CMS_MONGO_DB=people
CMS_MONGO_USERNAME=cms-service
CMS_MONGO_PASSWORD=<INSERT YOUR MONGODB PASSWORD>
CMS_MONGO_MAX_RETRY=5
CMS_MONGO_READ_TIMEOUT=5s
CMS_MONGO_WRITE_TIMEOUT=5s
CMS_MONGO_USE_TRANSACTION=false
CMS_REDIS_URL=localhost:6379
CMS_USE_CORS=true
CMS_DATA_MODE=mongo
CMS_CONTROLLER_PORT=8082
CMS_QUERY_WORKER_MODE=mongo
CMS_QUERY_WORKER_DB=events
CMS_QUERY_WORKER_USE_TRANSACTION=false
CMS_COMMAND_WORKER_MODE=mongo
CMS_COMMAND_WORKER_DB=events
CMS_COMMAND_WORKER_USE_TRANSACTION=false
```