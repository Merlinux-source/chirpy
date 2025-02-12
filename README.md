# Chirpy

A simple microblog backend learning project for Boot.dev.

## Quick Setup

This guide assumes you are using a Debian-based Linux distribution, but the general concepts apply to other systems,
including Windows.

1. **Update Your System**
   Ensure your system is up to date.

2. **Install Go**
   Install the Go build tools using either [the official docs](https://go.dev/doc/install)
   or [Webi](https://webinstall.dev/golang/).

3. **Clone the Repository**
   Clone the repository and navigate to the newly created directory:
   ```shell
   git clone https://github.com/Merlinux-source/chirpy
   cd chirpy
   ```

4. **Build the Go Executable**
   Compile the Go executable:
   ```shell
   go build .
   ```

5. **Set Up PostgreSQL**
   Set up a PostgreSQL server with a user and a database. [Click here for help](/docs/postgres-debian12.md).

6. **Create a `.env` File**
   Add a `.env` file with the following content:
   ```
   JWT_SECRET="RANDOMSTRINGHERE"
   DB_URL="postgres://postgres:postgres@localhost/chirpy?sslmode=disable"
   POKA_KEY="BootDevProvidedAPIKey"
   ```

7. **Update the `.env` File**
   Modify the entries in the `.env` file to match your environment. If you followed
   the [setup guide](/docs/postgres-debian12.md), you likely won't need to change the `DB_URL`.

8. **(Optional) Run Database Migrations**
   [Here is the guide for that](/docs/goose-migrations-debian12.md).
   > **Note:** This step is only necessary if you want to separate database permissions. The default settings are fine
   for development or review environments.

9. **(Optional) Move the Compiled Binary**
   Move the compiled binary to another location, such as `/opt/chirpy`:
   ```shell
   mkdir -p /opt/chirpy
   mv main /opt/chirpy/bin
   ```

10. **Run the Server**
    Execute the server binary:
    ```shell
    ./main # if you didn't move the binary
    /opt/chirpy/bin/main # if you did move the binary
    ```

Enjoy using Chirpy!