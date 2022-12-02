package echotron

import (
	"sync"
	"time"
)

type limit struct {
	counter int
	sended  bool
}

type Limiter struct {
	mutex    sync.Mutex
	LimitMap map[int64]*limit
	bot      API
}

func (l *Limiter) Increment(id int64) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if _, ok := l.LimitMap[id]; !ok {
		l.LimitMap[id] = &limit{
			counter: 0,
			sended:  false,
		}
	}

	l.LimitMap[id].counter++
}
func (l *Limiter) SetLimitTrue(id int64) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.LimitMap[id].sended = true
}

func (l *Limiter) GetCounter(id int64) int {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.LimitMap[id].counter
}

func (l *Limiter) GetSended(id int64) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.LimitMap[id].sended
}

func (l *Limiter) decrement(id int64) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	limit := l.LimitMap[id]

	if limit.counter > 0 {
		limit.counter--
	} else if limit.counter == 0 && limit.sended {
		limit.sended = false
		go l.bot.SendMessage("Możesz wysłać kolejną wiadomość.", id, nil)
	}

}

func (l *Limiter) Watcher() {
	for {
		for key := range l.LimitMap {
			l.decrement(key)
		}
		time.Sleep(2 * time.Second)
	}
}

func (l *Limiter) Check(id int64) bool {
	return l.GetCounter(id) < 5 && !l.GetSended(id)
}

func (l *Limiter) Init() {
	go l.Watcher()
}

func NewLimiter(bot API) *Limiter {
	return &Limiter{
		LimitMap: make(map[int64]*limit), bot: bot,
	}
}
