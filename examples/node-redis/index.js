const { exit } = require("process");
const { promisify } = require("util");
const redis = require("redis");

const client = redis.createClient({
    host: process.env.REDIS_URL
});

const setPair = promisify(client.set).bind(client);
const publish = promisify(client.publish).bind(client);

run();

async function run() {
    try {
        // Redis key/value store
        await setPair("ketap1", "Ketap");

        // Redis channel publish
        await publish("results", JSON.stringify({
            funcName: "func1", data: { ketap: "Ketap" }
        }));

        exit(0);
    } catch (e) {
        console.error(e);
        exit(1);
    }
}