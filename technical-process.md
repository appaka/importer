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
- postgres for transactions (or we could create an adhoc database, postgres is too much for only one table)
- mongodb for documents
- redis for caching queries?
- graphql for querying the documents?


---

MODULES:
- CSV reader
  - go app reading a CSV file and sending documents to MDM
  - delete not present documents?
    - send the PKs found in the CSV file to MDM, and MDM will delete the others
- MDM
  - receive documents and generate transactions
  - receive documents to be deleted, and generate transactions
  - transactions could trigger events
  - go app managing a MongoDB
- Cacher
  - receive transactions from a RabbitMQ queue (new/updated documents)
  - store documents into Redis
- GraphQL
  - Apollo access to the MDM MongoDB

