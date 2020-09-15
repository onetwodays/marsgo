import herolist from './herolist.json'

//整个文件需要 export 导出一个 js 对象。对象的 key 是由<Http_verb> <Resource_uri>
//构成的，值是 function，当一个 ajax 调用匹配了 key 后，与之对应的 function 就会被执行。函数中我们调用 res.json 就可以给浏览器返回结果
export default {
    '/apimock/web201605/js/herolist.json': [
        {
          ename: 106,
          cname: '小乔',
          title: '恋之微风',
          new_type: 0,
          hero_type: 2,
          skin_name: '恋之微风|万圣前夜|天鹅之梦|纯白花嫁|缤纷独角兽',
        },
      ],
      'POST /apimock/freeheros.json': (req, res) => {
        const { number } = req.body;
        function getRandomArrayElements(arr, count) {
          var shuffled = arr.slice(0),
            i = arr.length,
            min = i - count,
            temp,
            index;
          while (i-- > min) {
            index = Math.floor((i + 1) * Math.random());
            temp = shuffled[index];
            shuffled[index] = shuffled[i];
            shuffled[i] = temp;
          }
          return shuffled.slice(min);
        }
        const freeheros = getRandomArrayElements(herolist, number);
        res.send(freeheros);
      },

      'get /dev/random_joke': function (req, res) {
        res.status(500);
        res.json({});
      },

  };