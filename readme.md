# Mailing
___
## _test service for html-emails distribution_

### Features
- HTML-emails is filling by subscribers data (firstname, lastname, birthday)
- Emails can be sent now or with delay at the certain time
- Emails can be tracked for how many times and by whom it were opened by messageID
> Note: For `tracking` you need to host your server as it used [_the tracking pixel_](tracker/tracker.png) for that

### Tech
- Built in go version 1.18
- Uses [docker](https://www.docker.com/) container for _Postgres_ db with subscribers of distribution
- Uses [gin-gonic](https://github.com/gin-gonic/gin)
- Uses [pgx](https://github.com/jackc/pgx)


