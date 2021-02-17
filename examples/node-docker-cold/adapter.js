import { exit } from "process";

import { promisify } from "util";
import redis from "redis";

const client = redis.createClient({
    host: process.env.REDIS_URL
});

const publish = promisify(client.publish).bind(client);

export default async function triggerListener(func){
  var payload = process.env.PAYLOAD
  try{
    var obj = JSON.parse(payload)
  } catch(error){
    exit(1)
  }
  var resp = await func(obj) 

    // Redis channel publish
  await publish("results", JSON.stringify({funcName: process.env.FUNC_NAME, data: { ketap: resp }}));
  
  console.log(resp)

  exit(0)
}