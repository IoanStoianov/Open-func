import triggerListener from "./server.js";
// import http from "http";

import axios from 'axios';



const myFunc = async (args) => {
    //sample data validation
    if(args["Course"] == undefined){
        return "Invalid Input"
    }

    var resp = await axios.get('https://jsonplaceholder.typicode.com/todos/1')
            .then(response => { return response.data })
            .catch(error => {
                console.log(error);
            });
    return resp;
}

triggerListener(myFunc);