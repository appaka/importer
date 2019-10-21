# Glossary

- MDS: Master Data Store, where the master data is stored
- node/NDS: Node Data Store, exporters and importers machines, hosts sending/receiving data

# Importer process

## 1. input files

- define how to convert sources (xml, csv, json, services...) into documents
- validations: mandatory fields, available values, regexps...
- it could be a file or a service-call

## 2. document

- each row/line/piece of content of a import file generates a document
- a document is an entity which has to be sent to the system (product, category, order, page...)

## 3. diff

- the document is compared with the current data
- the generated diffs are sent to a queue to be processed and make the proper changes

## 4. queue

- the queue is processed
- tickets in queue have differents priorities (stock is higher, marketing data for products lower)

## 5. updater

- this is the queue worker
- get a piece to be changed, and execute the proper commands (insert/update/delete into tables)
- it sends the diff data to all nodes in the system


# Components

- sources
  - translations
  - validations
- entities definitions
- documents (maybe nosql objects)
- queues
- updater / worker
  - rules/tasks
  
# Technologies

- Go lang (this)
- RabbitMQ
- MongoDB
