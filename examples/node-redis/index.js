const redis = require("redis");

const client = redis.createClient({
    host: process.env.REDIS_URL
});

client.on("error", (err) => {
    console.error(err);
});

// Redis key/value store
client.set("ketap1", "Ketap", redis.print);

// Redis channel publish
client.publish("results", JSON.stringify({ ketap: "Ketap" }));