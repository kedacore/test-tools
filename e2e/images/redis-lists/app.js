const redis = require("redis");
const redisAddress = process.env.REDIS_ADDRESS
const redisHost = process.env.REDIS_HOST
const redisPort = process.env.REDIS_PORT != null ? parseInt(process.env.REDIS_PORT) : null
const listName = process.env.LIST_NAME
const redisPassword = process.env.REDIS_PASSWORD
const numberOfItemsToWrite = process.env.NO_LIST_ITEMS_TO_WRITE == null ? 100 : parseInt(process.env.NO_LIST_ITEMS_TO_WRITE)
const readProcessingTime = process.env.READ_PROCESS_TIME == null ? 1000 : parseInt(process.env.READ_PROCESS_TIME)


if (redisAddress == null && (redisHost == null || redisPort == null)) {
    throw new Error('Redis Address or Host and Port not set')
}

if (listName == null) {
    throw new Error('List name not set')
}

const redisConfig = {
    url: (redisAddress != null) ? `redis://${redisAddress}` : null,
    host: redisHost,
    port: redisPort,
    password: redisPassword
}

const client = createRedisClient(redisConfig)

var args = process.argv.slice(2);

if (args[0] == 'write') {
    writeToList(client, listName, numberOfItemsToWrite)
}
if (args[0] == 'read') {
    console.log(`reading from list '${listName}'`)
    readFromList(client, listName, readProcessingTime)
}

function createRedisClient(config) {
    const client = redis.createClient(config)

    if (config.password != null) {
        client.auth(config.password, (err, res) => {
            if (err != null) {
                console.error(`Authentication failed. Err is: ${err}`)
            }
        })
    }

    client.on("error", function (error) {
        console.error(error)
    });

    return client
}

async function writeToList(client, listName, numberOfItemsToWrite) {
    console.log(`writing to list '${listName}'`)
    for (var i = 0; i < numberOfItemsToWrite; i++) {
        await writeItemToList(client, listName, i.toString())
    }
    console.log(`done writing to list '${listName}'`)
    client.quit((err, res) => {
        if (err != null) {
            console.error(err)
        }
    })
    process.exit()
}

async function writeItemToList(client, listName, item) {
    return new Promise((resolve, reject) => {
        client.lpush([listName, item], (err, res) => {
            if (err != null) {
                console.error(err)
                reject()
            } else {
                console.log(`added item ${item}`)
                resolve()
            }
        })
    })
}

function readFromList(client, listName, readProcessingTime) {
    client.llen(listName, (err, reply) => {
        console.log(`list size: ${reply}`)
        if (reply > 0) {
            client.lpop(listName, (err, reply) => {
                console.log(`read item: ${reply}`)
            });
        }
        setTimeout(function () {
            readFromList(client, listName, readProcessingTime)
        }, readProcessingTime);

    })
}