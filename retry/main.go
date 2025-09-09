package retry

import (
	"log"
	"time"
)

func ForeverSleep(fn func() error) {
	for {
		if err := fn(); err != nil {
			log.Println("retrying to connect with database...")
			time.Sleep(time.Second*2)
			continue
		}
		return	
	}
}