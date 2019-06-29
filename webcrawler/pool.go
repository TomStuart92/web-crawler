package webcrawler

import (
	"io"
	"log"
)

type Pool struct {
	concurency  int
	workChan    chan string
	resultsChan chan PoolJob
	errChan     chan error
	workFn      func(string) (error, io.Reader)
	transformFn func(io.Reader) []string
}

type PoolJob struct {
	Input  string
	Output []string
}

func InitializePool(concurency int) *Pool {
	log.Printf("Initialized Pool with %d workers", concurency)
	return &Pool{concurency, nil, nil, nil, nil, nil}
}

func (p *Pool) SetWorkFn(workFn func(string) (error, io.Reader)) {
	p.workFn = workFn
}

func (p *Pool) SetTransformFn(transformFn func(io.Reader) []string) {
	p.transformFn = transformFn
}

func (p *Pool) Process(workChan chan string) (chan PoolJob, chan error) {
	p.workChan = workChan
	p.resultsChan = make(chan PoolJob)
	p.errChan = make(chan error)

	for i := 0; i < p.concurency; i++ {
		go func(workerID int) {
			for link := range p.workChan {
				log.Printf("Worker %d, processing %s", workerID, link)
				err, result := p.workFn(link)
				if err != nil {
					go func() { p.errChan <- err }()
					continue
				}
				go func(link string) {
					transformed := p.transformFn(result)
					p.resultsChan <- PoolJob{Input: link, Output: transformed}
				}(link)
			}
		}(i)
	}
	return p.resultsChan, p.errChan
}

func (p *Pool) Stop() {
	if p.workChan == nil {
		log.Print("No Jobs In Progress")
		return
	}
	log.Print("Gracefully Terminating Pool...")
	close(p.workChan)
	close(p.resultsChan)
	close(p.errChan)
}
