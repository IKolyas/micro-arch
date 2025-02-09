package actions

import (
	"sync"
	"time"

	"github.com/IKolyas/otus-highload/internal/domain"
	"github.com/go-faker/faker/v4"
)

func CreateRandomUsers(count int, repo domain.Repository[domain.User]) (success int, errors int) {
	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	workers := 10 // количество горутин, можно настроить по необходимости
	usersPerWorker := count / workers

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			localSuccess := 0
			localErrors := 0

			for range usersPerWorker {
				date, err := time.Parse("2006-01-02", faker.Date())
				if err != nil {
					panic(err)
				}
				user := domain.User{
					Login:      faker.Email(),
					Password:   faker.Password(),
					FirstName:  faker.FirstName(),
					SecondName: faker.LastName(),
					Gender:     "Y",
					Birthdate:  date,
					Biography:  faker.Paragraph(),
					City:       faker.Word(),
				}

				_, err = repo.Save(&user)
				if err != nil {
					localErrors++
				} else {
					localSuccess++
				}
			}

			mu.Lock()
			success += localSuccess
			errors += localErrors
			mu.Unlock()
		}(w) // передаем идентификатор горутины, если нужно
	}

	wg.Wait()
	return success, errors
}
