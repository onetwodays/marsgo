//约定 mock 目录里所有的 .js 文件会被解析为 mock 文件,后在浏览器里访问 http://127.0.0.1:8000/api/users 就可以看到 ['a', 'b'] 了.
// 整个文件需要 export 出一个js 对象。对象的 key 是由 <Http_verb> <Resource_uri>;值是 function，当一个 ajax 调用匹配了 key 后，与之对应的 function 就会被执行。函
export default {
    '/api/users': ['a', 'b'],
  }

  //export default {
  //  'get /dev/random_joke': function (req, res) {
  //    const responseObj = random_jokes[random_joke_call_count % random_jokes.length];
  //    random_joke_call_count += 1;
  //    setTimeout(() => {
  //      res.json(responseObj);
  //    }, 3000);
  //  },
  //};