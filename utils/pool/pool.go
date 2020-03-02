package pool

import (
    "io"
    "log"
    "sync"
    "errors"
)

var(
    ErrPoolClosed = errors.New("资源池已经关闭。")
)

//一个安全的资源池，被管理的资源必须都实现io.Close接口
type Pool struct {
    m sync.Mutex // 互斥锁,这主要是用来保证在多个goroutine访问资源时，池内的值是安全的
    res chan io.Closer // 一个有缓冲的通道，用来保存共享的资源，这个通道的大小，在初始化Pool的时候就指定的。注意这个通道的类型是io.Closer接口，所以实现了这个io.Closer接口的类型都可以作为资源，交给我们的资源池管理
    factory func() (io.Closer,error) //一个函数类型，它的作用就是当需要一个新的资源时，可以通过这个函数创建，也就是说它是生成新资源的，至于如何生成、生成什么资源，是由使用者决定的，所以这也是这个资源池灵活的设计的地
    closed bool  //表示资源池是否被关闭，如果被关闭的话，再访问是会有错误的。
}


//创建一个资源池
func NewPool(fn func() (io.Closer, error), size uint) (*Pool, error) {
    if size <= 0 {
        return nil, errors.New("size的值太小了。")
    }
    return &Pool{
        factory: fn,
        res:     make(chan io.Closer, size),
    }, nil
}


//从资源池里获取一个资源 select的多路复用，因为这个函数不能阻塞，可以获取到就获取，不能就生成一个
func (p *Pool) Acquire() (io.Closer,error) {
    select {
    case r,ok := <-p.res:
        log.Println("Acquire:共享资源")
        if !ok {
            return nil,ErrPoolClosed
        }
        return r,nil
    default:  // 如果没有资源，则调用factory方法生成一个并返回
        log.Println("Acquire:新生成资源")
        return p.factory()
    }
}


//关闭资源池，释放资源  互斥锁，因为有个标记资源池是否关闭的字段closed需要再多个goroutine操作，
//所以我们必须保证这个字段的同步。这里把关闭标志置为true
//释放通道中的资源，因为所有资源都实现了io.Closer接口，所以我们直接调用Close方法释放资源即可
func (p *Pool) Close() {
    p.m.Lock()
    defer p.m.Unlock()

    if p.closed {
        return
    }

    p.closed = true

    //关闭通道，不让写入了
    close(p.res)

    //关闭通道里的资源
    for r:=range p.res {
        r.Close()
    }
}

/*
释放资源本质上就会把资源再发送到缓冲通道中，就是这么简单，
不过为了更安全的实现这个方法，我们使用了互斥锁，保证closed标志的安全，
而且这个互斥锁还有一个好处，就是不会往一个已经关闭的通道发送资源。

这是为什么呢？因为Close和Release这两个方法是互斥的，Close方法里对closed标志的修改，
Release方法可以感知到，所以就直接return了，不会执行下面的select代码了，
也就不会往一个已经关闭的通道里发送资源了。

如果资源池没有被关闭，则继续尝试往资源通道发送资源，如果可以发送，就等于资源又回到资源池里了；
如果发送不了，说明资源池满了，该资源就无法重新回到资源池里，那么我们就把这个需要释放的资源关闭，抛弃了
 */
func (p *Pool) Release(r io.Closer){
    //保证该操作和Close方法的操作是安全的
    p.m.Lock()
    defer p.m.Unlock()

    //资源池都关闭了，就省这一个没有释放的资源了，释放即可
    if p.closed {
        r.Close()
        return
    }

    select {
    case p.res <- r:
        log.Println("资源释放到池子里了")
    default:
        log.Println("资源池满了，释放这个资源吧")
        r.Close()
    }
}

/*
func main() {
	//等待任务完成
	var wg sync.WaitGroup
	wg.Add(maxGoroutine)

	p, err := common.New(createConnection, poolRes)
	if err != nil {
		log.Println(err)
		return
	}
	//模拟好几个goroutine同时使用资源池查询数据
	for query := 0; query < maxGoroutine; query++ {
		go func(q int) {
			dbQuery(q, p)
			wg.Done()
		}(query)
	}

	wg.Wait()
	log.Println("开始关闭资源池")
	p.Close()
}

//模拟数据库查询
func dbQuery(query int, pool *common.Pool) {
	conn, err := pool.Acquire()
	if err != nil {
		log.Println(err)
		return
	}

	defer pool.Release(conn)

	//模拟查询
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	log.Printf("第%d个查询，使用的是ID为%d的数据库连接", query, conn.(*dbConnection).ID)
}
//数据库连接
type dbConnection struct {
	ID int32//连接的标志
}

//实现io.Closer接口
func (db *dbConnection) Close() error {
	log.Println("关闭连接", db.ID)
	return nil
}

var idCounter int32

//生成数据库连接的方法，以供资源池使用
func createConnection() (io.Closer, error) {
	//并发安全，给数据库连接生成唯一标志
	id := atomic.AddInt32(&idCounter, 1)
	return &dbConnection{id}, nil
 */






