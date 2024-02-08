**WARNING: This is an alpha release - don't use it for essential files. Suggestions are welcome!**

# What is this program for?

The goal of this project is to provide a simple way to create an incremental backup and upload
it to some cloud storage.

## How `glacier_backup` works?

The goal is to create a remote backup of local files. To achive this it will execute following steps:

1. Check if previous backup, for given source was successfull and if can proceed with creating new one,
2. Create a local, increamental `tar` archive,
3. Encrypt local backup,
4. Upload archive to `AWS Glacier`,
5. Remove local copy,
6. Save progress to local database.

## How this software checks if backup is not corrupted?

1. The status of previous backups is stored in the local database. In case the previous backup was not successful/rollbacked, `glacier_backup` will prevent the creation of new backups.
2. If any of the steps fail, we will try to rollback all changes to leave everything in the last save state.
3. When uploading data to `AWS glacier`, we check data integrity by comparing the checksum of the local copy with that received from `AWS`.

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
      canChange: true
      keep: false
      vault: "glacierVaultName"
      password: "1234"
```

The backup section is responsible for holding all backup configurations. Every entry consists of the following elements:
* src - source directory or file you wish to create a backup for,
* dst - destination path (must exist) where the temporary backup archive will be stored before uploading it to remote storage,
* canChange - specifies if files in the source directory can change while creating a backup - if set to false (default value), the backup will fail if any file changes,
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

## Local db configuration

`glacier_backup` stores some important metadata in the local database. This data is crucial for the correct operation of the program. In the DB are stored:

* status of workflows,
* status of jobs in given workflow,
* time of workflow creation.

This data will be used to determine if the previous backup workflow was successful and if the new one can be safely created.

Example configuration:

```yaml
db:
  path: "/path/to/local/db"
```

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

**Very important!**

Creating a backup will result in creating a file with a `.manifest` extension. It is essential not to delete it. It is used by `tar` to create an incremental backup. If you delete it, the next backup job will create a tarball with all files from the source directory.
