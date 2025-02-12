# Manual Database Migration

This step is necessary only if you provided a connection string with a user that does not own the database. For official
documentation, refer to [here](https://pressly.github.io/goose/installation/#linux), but note that my version may
deviate from best practices.

1. **Download the Installer**
   Navigate to a temporary directory and download the installer:
   ```shell
   cd /tmp
   mkdir goose_install && cd goose_install
   curl -fsSLO https://raw.githubusercontent.com/pressly/goose/master/install.sh
   chmod u+x install.sh
   ```

2. **Set Installation Directory (Optional)**
   If you need to change the installation directory, set the environment variable `GOOSE_INSTALL`. I recommend auditing
   the script, as it is short and easy to understand.

   > **Note:**  
   > Auditing any install script is crucial, as even trusted software can be compromised. A thorough review helps
   protect your system from hidden vulnerabilities or malicious code.

3. **Run the Installer**
   Execute the installer to complete the installation:
   ```shell
   ./install.sh
   ```

4. **Start Migrations**
   Navigate back to your cloned project and set the necessary environment variables:
   ```shell
   # Update GOOSE_DBSTRING according to your environment.
   export GOOSE_DBSTRING='postgres://postgres:postgres@127.0.0.1:5432/chirpy'
   export GOOSE_DRIVER=postgres
   cd sql/schema
   goose up
   ```

If there are no errors, you can proceed with the setup.