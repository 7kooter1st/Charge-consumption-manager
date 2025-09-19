package repository

import (
  "fmt"
  "strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
  return &Repository{}, nil
}

type Usecase struct { // вот наша новая структура 
  ID    int // поля структур, которые передаются в шаблон
  Title string // ОБЯЗАТЕЛЬНО должны быть написаны с заглавной буквы (то есть публичными)
  Adres string
  Exprnditure uint
}

func (r *Repository) GetUsecase(id int) (Usecase, error) {
	// тут у вас будет логика получения нужной услуги, тоже наверное через цикл в первой лабе, и через запрос к БД начиная со второй 
	usecases, err := r.GetUsecases()
	if err != nil {
		return Usecase{}, err // тут у нас уже есть кастомная ошибка из нашего метода, поэтому мы можем просто вернуть ее
	}

	for _, usecase := range usecases {
		if usecase.ID == id {
			return usecase, nil // если нашли, то просто возвращаем найденный заказ (услугу) без ошибок
		}
	}
	return Usecase{}, fmt.Errorf("заказ не найден") // тут нужна кастомная ошибка, чтобы понимать на каком этапе возникла ошибка и что произошло
}

func (r *Repository) GetUsecases() ([]Usecase, error) {
  // имитируем работу с БД. Типа мы выполнили sql запрос и получили эти строки из БД
  usecases := []Usecase{ // массив элементов из наших структур
    {
      ID:    1,
      Title: "игры",
      Adres: "http://127.0.0.1:9000/lab1/premium_photo-1682125220008-006e43a6896a.jpeg",
      Exprnditure: 100,
    },
    {
      ID:    2,
      Title: "просмотр видео",
      Adres: "http://127.0.0.1:9000/lab1/video.jpeg",
      Exprnditure: 20,
    },
    {
      ID:    3,
      Title: "просмотр reels",
      Adres: "http://127.0.0.1:9000/lab1/instagram.jpeg",
      Exprnditure: 22,
    },
    {
      ID:    4,
      Title: "мессенджеры",
      Adres: "http://127.0.0.1:9000/lab1/messanger.jpeg",
      Exprnditure: 20,
    },
    {
      ID:    5,
      Title: "навигатор", 
      Adres: "http://127.0.0.1:9000/lab1/maps.jpeg",
      Exprnditure: 1448,
    },
    {
      ID:    6,
      Title: "звонки",  
      Adres: "http://127.19.0.1:9000/lab1/phone.jpeg",
      Exprnditure: 427,
      },    
    {
      ID:    7,
      Title: "видеозвонки",
      Adres: "http://127.19.0.1:9000/lab1/zoom.jpeg",
      Exprnditure: 100,
    },
    {
      ID:    8,
      Title: "музыка",  
      Adres: "http://127.19.0.1:9000/lab1/music.jpeg",
      Exprnditure: 15,
    },
  }
  // обязательно проверяем ошибки, и если они появились - передаем выше, то есть хендлеру
  // тут я снова искусственно обработаю "ошибку" чисто чтобы показать вам как их передавать выше
  if len(usecases) == 0 {
    return nil, fmt.Errorf("массив пустой")
  }

  return usecases, nil
}

func (r *Repository) GetUsecasesByTitle(title string) ([]Usecase, error) {
	usecases, err := r.GetUsecases()
	if err != nil {
		return []Usecase{}, err
	}

	var result []Usecase
	for _, usecase := range usecases {
		if strings.Contains(strings.ToLower(usecase.Title), strings.ToLower(title)) {
			result = append(result, usecase)
		}
	}

	return result, nil
}