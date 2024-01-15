# What is this program for?

The goal of this project is to provide a simple way to create an incremental backup and upload
it to some cloud storage.

# Current functionalities

Currently, you can:

* Create an incremental backup,
* Encrypt backup,
* Upload backup to the AWS Glacier.

# Dependencies

This program requires the following programs to run:

* GNU tar - https://savannah.gnu.org/git/?group=tar
* GNU Privacy Guard - https://www.gnupg.org/

If you are using `GNU/Linux` distribution with `apt` package manager, you can (probably) install both using:

```bash
> sudo apt install tar
> sudo apt install gnupg
```

# Configuration

You will be asked to provide a path to the configuration file when running this program.
You can find the example configuration in `example/config.`

## Backup configuration

An example backup configuration:

```yaml
backup:
    - src: "/path/to/backup/src"
      dst: "/path/to/backup/dst"
      keep: false
      vault: "glacierVaultName"
      password: "1234"
```

The backup section is responsible for holding all backup configurations. Every entry consists of the following elements:
* src - source directory or file you wish to create a backup for,
* dst - destination path (must exist) where the temporary backup archive will be stored before uploading it to remote storage,
* keep - specifies if you wish to keep the local copy of the backup
* vault - the name of AWS Glacier Vault
* password - password used for backup encryption

## AWS configuration

An example of AWS configuration:

```yaml
aws:
  profile: "aws_profile_name"
  account: "123456789"
```

There are only two settings: profile and account. The `profile` specifies the profile you configured while setting up AWS, while the `account` is your AWS account.

# Running glacier_backup:

First, download the latest release. I suggest adding it to `/usr/bin` so it will be available without providing the absolute path:

```bash
> sudo cp /path/to/downloaded/binary/glacier_backup /usr/bin
> sudo chmod +x /usr/bin/glacier_backup
```

Now, you should be able to run your backup jobs with the following:

```bash
> glacier_backup /path/to/config/config.yml
```
