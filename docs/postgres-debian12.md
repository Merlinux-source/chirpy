# Setting Up PostgreSQL for This Project

Log in as the root user or execute `sudo su -`. This guide assumes you have root privileges.

1. **Update Your System and Install PostgreSQL**
   Ensure your system is updated and install the server software:
   ```shell
   apt update && apt upgrade -y
   apt install postgresql-common -y
   ```

2. **Set Up the Database User**
   Change the password for the PostgreSQL user:
   ```shell
   passwd postgres
   ```
   Enter a password when prompted. **Note:** The password will not be displayed on the screen, which is normal. I
   recommend using `postgres` as the password for simplicity since the default configuration of postgres restricts
   connections to the localhost anyway.

3. **Create the Database**
   Switch to the PostgreSQL user and create the database:
   ```shell
   su postgres
   psql
   ```
   In the PostgreSQL prompt, run:
   ```postgresql
   ALTER USER postgres WITH PASSWORD 'postgres';
   CREATE DATABASE chirpy OWNER postgres;
   \q
   ```

You now have a working PostgreSQL installation.