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

# transactionalization / t18n

We receive data (documents, rows, items...) and generate transactions in the timeline:
- what: what process created this transaction (import file asdf.csv 19.05.2019, tcp connection from qwerty.com...)
- when: datetime
- who: nodes and users
- where: document type and ID (product:123, category:69, order:102345...)
- data: field and value of the new version (title="asasdsadasd", price="123", ...)
- last_version: reference to the transaction with the previous version (null if new)

We need:
- postgres for transactions
- mongodb for documents
- redis for caching
- graphql for querying?
