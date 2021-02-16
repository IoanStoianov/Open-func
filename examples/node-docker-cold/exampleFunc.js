import triggerListener from "./adapter.js";

import axios from 'axios';

import { promisify } from "util";
import redis from "redis";

const client = redis.createClient({
    host: process.env.REDIS_URL
});

const publish = promisify(client.publish).bind(client);

const myFunc = async (args) => {
    //sample data validation
    // if(args["Course"] == undefined){
    //     return "Invalid Input"
    // }


    // Redis channel publish
    await publish("results", JSON.stringify(args));

    var resp = await axios.get('https://jsonplaceholder.typicode.com/todos/1')
            .then(response => { return response.data })
            .catch(error => {
                console.log(error);
            });

    return resp;
}

triggerListener(myFunc);