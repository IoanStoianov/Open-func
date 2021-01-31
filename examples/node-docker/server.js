var express = require('express');
var bodyParser = require('body-parser');

var app = express();


app.use(bodyParser.json());

app.post('/triggerHttp', (req, res) => {
    console.log('Got body:', req.body);
    res.send(["Ketap",req.body]);
});

app.listen(process.env.OPEN_FUNC_PORT || 3014, () => {
  console.log('Listening to Port ' + process.env.OPEN_FUNC_PORT || 3014);
});