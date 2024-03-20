# Использование 🛠️

![cat_Persik_033](https://github.com/shuklarituparn/Conversion-Microservice/assets/66947051/6c1ab754-da77-4133-b3b8-b2554f3cd361)


`Если вы тестируете api на swagger, в токене авторизации, пожалуйста, введите "Bearer <ваш токен>".`

---

## Вход:

Существует два типа пользователей:


- Админ
- Обычный пользователь


Для администратора адрес электронной почты и пароль по умолчанию следующие:

`email@example.com` и `adminpassword`

Администратор и пользователи могут войти в систему, отправив post-запрос на конечную точку: `https://api.rtprnshukla.ru/api/v1/users/login`


измените `api.rtprnshukla.ru` на адрес `localhost` при запуске локально

``` json 
{
  "email": "admin@example.com",
  "password": "adminpassword"
}

```
---
### Регистрвция нового пользователя

Пользователь может зарегистрироваться, отправ post-запрос в конечную точку со следующим телом json

`https://api.rtprnshukla.ru/api/v1/users/register`

```json

{
  "email": "user12@example.com",
  "password": "userpassword"
}

```
---

Роль пользователя по умолчанию после регистрации:  "USER"



## Актеры:


### Создайте актера (POST)


Админ может создать актера, отправ POST запрос к конечной точке со следующим телом json после авторизации


`https://api.rtprnshukla.ru/api/v1/actors/create`


```json

{
  "first_name": "New Actor",
  "last_name": "Surname",
  "gender": "male",
  "birth_date": "1969-03-11"
}


```
---
### обновите актера (UPDATE)

Админ может обновлять актера, отправ UPDATE запрос к конечной точке со следующим телом json после авторизации

`https://api.rtprnshukla.ru/api/v1/actors/update`

```json

{
  "id": 8,
  "first_name": "New Actor",
  "last_name": "Surname",
  "gender": "male",
  "Description": "min",
  "Genre": "required",
  "Rating":"1.0" ,
  "ReleaseYear": "2000",
  "Title": "A new movie",
  "birth_date": "1969-03-11",
  "movies": [{
  "id": 2,
  "title": "Fight Club",
  "release_year": 2000,
  "genre": "Drama",
  "description": "An insomniac office worker and a devil-may-care soapmaker form an underground fight club that evolves into something much, much more.",
  "rating": 10
  }]
}


```
---

### Удалите актера (DELETE)


Админ может удалять актера, отправ DELETE запрос к конечной точке со следующим телом json после авторизации


`https://api.rtprnshukla.ru/api/v1/actors/delete?id={id актера}`

Необходимо ввести идентификатор актера после слова id

---

### Получите актера (GET)

Админ/Пользователь может получите всех актеров, отправ GET запрос к конечной точке


`https://api.rtprnshukla.ru/api/v1/actors/get?id={id актера}`


Необходимо ввести идентификатор актера после слова id

---


### получите всех актеров (GET)

Админ/Пользователь может получит актера, отправ GET запрос к конечной точке

`https://api.rtprnshukla.ru/api/v1/actors/all?page_size={размер}&page={Номер страницы}&sort_by=first_name&sort_order=ASC`

Поддерживает разбиение на страницы и сортировку


Поле для сортировки (по умолчанию birth_date)

Порядок сортировки (ASC или DESC, по умолчанию DESC)


---

### Исправьте актера (PATCH)


Админ может исправить актера, отправ PATCH запрос к конечной точке со следующим телом json после авторизации

`https://api.rtprnshukla.ru/api/v1/actors/patch?id={id актера}`


Необходимо ввести идентификатор актера после слова id


```json

{
  "first_name": "New Actor Rename",
  "movies": [{"id":10}]
}


```

---




## Фильмы:


### Создайте Фильм (POST)


Админ может создать Фильм, отправ POST запрос к конечной точке со следующим телом json после авторизации


`https://api.rtprnshukla.ru/api/v1/movies/create`


```json
{
 "title": "New Movie",
 "release_year": 2000,
 "genre": "Drama",
 "description":"Description Movie",
 "rating": 8.9
}


```
---
### обновите Фильм (UPDATE)

Админ может обновлять Фильм, отправ UPDATE запрос к конечной точке со следующим телом json после авторизации

`https://api.rtprnshukla.ru/api/v1/movies/update`

```json

{
  "id": 8,
  "title": "Fight Club",
  "release_year": 1999,
  "genre": "Drama",
  "description": "An insomniac office worker and a devil-may-care soapmaker form an underground fight club that evolves into something much, much more.",
  "rating": 8.9,
  "actors": [{
    "id": 8,
    "birth_date": "1937-06-01",  
    "first_name": "John",          
    "gender": "Male",
    "last_name": "Crews"
  }]
}


```
---

### Удалите Фильм (DELETE)


Админ может удалить Фильм, отправ DELETE запрос к конечной точке со следующим телом json после авторизации


`https://api.rtprnshukla.ru/api/v1/movies/delete?id={id актера}`

---

### Получите Фильм (GET)

Админ/Пользователь может получит Фильм, отправ GET запрос к конечной точке


`https://api.rtprnshukla.ru/api/v1/movies/get?id={id актера}`


Необходимо ввести идентификатор актера после слова id

---


### получите все Фильмы (GET)

Админ/Пользователь может получит все Фильмы, отправ GET запрос к конечной точке

`https://api.rtprnshukla.ru/api/v1/movies/all?page_size={размер}&page={Номер страницы}&sort_by=first_name&sort_order=ASC`

Поддерживает разбиение на страницы и сортировку


Поле для сортировки (по умолчанию рейтинг)

Порядок сортировки (ASC или DESC, по умолчанию DESC)



---

### Исправьте Фильм (PATCH)


Админ может исправить Фильм, отправ PATCH запрос к конечной точке со следующим телом json после авторизации

`https://api.rtprnshukla.ru/api/v1/movies/patch?id={id актера}`


Необходимо ввести идентификатор актера после слова id


```json
{
 "description":"An insomniac office worker and a devil-may-care soapmaker form an underground fight club that evolves into something much, much more.",
 "rating": 9.0,
 "actors": [{"id": 5}]
}
```

---


## Вклад 🤝

Ваш вклад приветствуется!

## Лицензия 📄

Этот проект лицензирован под [лицензией MIT](LICENSE).
