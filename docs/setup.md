# Установка ⚙️

![944960 512 (1)](https://github.com/shuklarituparn/Conversion-Microservice/assets/66947051/e31ed3cb-cfa1-454a-b664-5a2e63c579e3)


## Сайт  🌐


API и все методы доступны здесь: 

`https://api.rtprnshukla.ru/`


Docs доступны здесь:

`https://api.rtprnshukla.ru/docs`
`

pgAdmin доступен здесь:

`https://pgadmin.rtprnshukla.ru/docs`


по умолчанию id и пароль графаны:

```
admin@example.com
adminpassword
```




метрики доступны здесь:

`https://monitor.rtprnshukla.ru/metrics`

grafana доступна здесь:


`https://monitor.rtprnshukla.ru/`


по умолчанию id и пароль графаны:  

```
admin
T=5,Wg,epq5;vfU
```





## docker-compose 🚀

* Клонируйте проект, выполнив следующую команду:

  `git clone git@github.com:shuklarituparn/Filmoteka.git`


* Теперь выполните следующую команду:  `cd Filmoteka`



* Чтобы установить все зависимости, запустите следующую команду, находясь в корневой директории проекта:

  ` make setup `


* Заполните файл `.env.example` и переименуйте его в `.env` (Как заполнить env, я пишу ниже)

* Заполните следующее в файле docker compose

      В сервисе grafana заполните следующие поля
      
      - GF_SMTP_ENABLED=true
      - GF_SMTP_HOST=< host-smtp:порт (например smtp.resend.com:587) >
      - GF_SMTP_USER=< Ваш юзернейм/email от почтового провайдера>
      - GF_SMTP_PASSWORD=< Ваш пароль >
      - GF_SMTP_SKIP_VERIFY=false
      - GF_SMTP_FROM_NAME=Grafana
      - GF_SMTP_FROM_ADDRESS=<Ваш юзернейм/email от почтового провайдера>
      
      В сервисе postgres заполните следующее
      
      - POSTGRES_USER: <Юзернэм вашего постгреса>
      - POSTGRES_PASSWORD: <Пароль вашего постгреса>
      - POSTGRES_DB: <Название ваши базы данных>
      
      (ваше юзернэм, пароль вашего postgres и название вашей базы данных должны совпадать с вашими переменными в файле .env )
      
      В сервисе pgAdmin заполните следующее
      
      - PGADMIN_DEFAULT_EMAIL: <Адрес электронной почты по умолчанию для pgAdmin>
      - PGADMIN_DEFAULT_PASSWORD: <Пароль по умолчанию для pgAdmin>

* Находясь в в корневой директории проекта, выполните следующую команду, чтобы запустить: `docker compose up`

  `  Убедитесь, что у вас установлен Docker перед выполнением вышеуказанной команды! `


* Cервис доступен по адресу `http://localhost:8085`

  `(Без  метода http://localhost:8085 показывает проверку работоспособности сервиса)`

* Вы можете получить доступ к метрикам prometheus по адресу `localhost:8085/metrics`

* Графана доступна по адресу `localhost:3000`

* pgAdmin доступен по адресу `http://localhost:16543`

`При регистрации сервера в pgAdmin use введите имя сервиса вашего postgres из docker-compose в качестве имени вашего host `


---



## ENV Файл 📝

      - POSTGRES_USER: <Юзернэм вашего постгреса>
      - POSTGRES_PASSWORD: <Пароль вашего постгреса>
      - POSTGRES_DB: <Название ваши базы данных>
      - POSTGRES_HOST: <введите имя сервиса вашего postgres из docker-compose>
      - POSTGRES_PORT=5432 (по умолчанию)
      
      (ваше юзернэм, пароль вашего postgres и название вашей базы данных должны совпадать с вашими переменными в файле docker-compose )
      
      
     
        - JWT_SECRET= < чтобы подписать токены> 
        
        (Например: SGVsbG8sIEhvdyBhcmUgeW91LCBNeSBuYW1lIGlzIFJpdHVwYXJuLiBJIGFtIGEgc3R1ZGVudCBpbiBNSVBUIGFuZCBBbWJhc3NhZG9yIG9mIFZLIGNvbXBhbnk)

        - PORT= 8085 (По умолчанию, рекомендуется оставить его как есть, если вы не хотите изменить это везде)
        
