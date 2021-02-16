import { exit } from "process";

export default async function triggerListener(func){
  var payload = process.env.PAYLOAD
  try{
    var obj = JSON.parse(payload)
  } catch(error){
    exit(1)
  }
  var resp = await func(obj) 
  console.log(resp)

  exit(0)
}