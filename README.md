# Glossary

- MDS: Master Data Store, where the master data is stored
- node/NDS: Node Data Store, exporters and importers machines, hosts sending/receiving data

# Importer process

## 1. MDS receives data

- Any node can connect to MDS to send data (csv file, API Rest...).
- Any sent data is related to a type (product, category, stock...).
- MDS process the data, detects differences...
- If there is any change in the data (diff) then MDS sends data to the listeners
- The data could be the full row/document, just some fields, calculated fields, the diff...
- And the format of the sent data could be a csv file, a call to any API, etc
- If a node wants to receive new changes it should register as listener by the proper type of object

### 1.1. Receive files

- configured by a config file
- local directories for sources, archive, fails, success...
- remotes: FTP, SSH, NFS...
- regexp to match the files
- type: csv, tsv, Excel, json, yml...
- indexes fields
- header field names translations (field name in file is related to another field name in MDS)
- cron times (how often it should check for new files)
- copying files check (check if the file is uploading)
- validation rules: mandatory fields, regexps for values...

### 1.2. Receive data by API

- TODO

### 1.3. Receive data from a RabittMQ channel

- TODO

## 2. Node receives data

- A node subscribed to a type of object (product, category, stock...) will receive any change about any object of that type
- A node could be subscribed just to any objects of a type (product 543, category 987, ...)
- The data could be received as a diff, as the full row, just a few fields (configured by config file)
- A node could be a channel of RabbitMQ

# Configuration files

- master.conf : MDS configuration, data received
- nodes/999-node-name.conf : node configuration, data received


# How to receive the data from the MDS

When the MDS receives data and detect changes, those changes are sent inmediatly to their subscribers.

The subscribers should choose the way this data is received:
- just the diff / all the row / certain fields of the row / transformed fields
- by json, csv, yml...
- remotes : FTP, SSH, NFS...

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

# Deprecated

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

