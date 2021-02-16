import express from 'express';

export default async function triggerListener(func){
  const app = express();

  app.use(express.json());

  app.post('/triggerHttp', async (req, res) => {
    console.log('Passing body to function:', req.body);
    let response = await func(req.body)
    // TODO log response in table
    res.send(response)
  });

  app.listen(process.env.OPEN_FUNC_PORT || 3014, () => {
    console.log('Listening on port ' + (process.env.OPEN_FUNC_PORT || 3014));
  });
}