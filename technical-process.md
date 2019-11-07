# node server

## data in
- listening for new data coming (own TCP/IP, or look at the RabbitMQ protocol)
- watchers in local filesystem to import files
- check periodically remote filesystems (FTP, NFS...)

## data process
- read from different files formats
- make translations/transformations
- generate diffs
- save data into db

## data out
- generate files with all diffs and all modified rows
- send files to local and remote filesystems (FTP, NFS...)
- send diffs/rows by TCP to all the document subscribers
