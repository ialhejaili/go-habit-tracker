# go-habit-tracker

A Go-based command-line application designed to help users manage their daily habits. With user
authentication and a PostgreSQL database, it allows users to register, log in, add, view,
delete, and mark habits as completed for the day. The app features a "Today's Must Do" option, showing
incomplete habits for the current day.

### Set Environment Variables
Before running the application, you need to create a `.env` file in the root directory of the project. This file should contain your database connection strings for both the production and test environments.

1- Create a .env file in the root directory of your project:
```sh
   touch .env
```

2- Add the following variables to your .env file:
```
   PG_DSN=production_db_connection
   TEST_PG_DSN=test_db_connection
```

Replace production_db_connection and test_db_connection with your actual database connection strings.
