const express = require('express');

const app = express();

app.use(express.json());

app.get('/', (_, res) => {
  res.send('Ketap works...');
});

app.post('/triggerHttp', (req, res) => {
  console.log('Got body:', req.body);
  res.send(["Ketap", req.body]);
});

app.listen(process.env.OPEN_FUNC_PORT || 3014, () => {
  console.log('Listening on port ' + process.env.OPEN_FUNC_PORT || 3014);
});