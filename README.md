# Самок-Аах-т

Приложение для заказа стриптизеров и стриптизерш в любую точку России.

## Функциональные требования

1. **Регистрация клиента/модели**
2. **Аккаунт клиента:** 
    - Имя
    - Фото по желанию
    - Дата рождения
    - Контактные данные (телефон/почта + паспорт (жестко 18+))
    - История заказов 
    - Реферальная система:
       - Уникальный реферальный код
       - Баланс - вывод денег невозможен
       - Пополнение баланса (ввод реквизитов - в бд не нужно)
       - История приглашений
3. **Аккаунт модели:**
    - Псевдоним
    - Фото обязательно
    - Возраст
    - Контактные данные (телефон/почта/соц сети мб + паспорт (контроль))
    - Рейтинг
    - Отзывы
    - Портфолио (по желанию)
    - Список услуг + цена
    - Реферальная система:
       - Уникальный реферальный код
       - Баланс
       - Пополнение баланса
       - История приглашений
4. **Администратор:**
    - Управление пользователями (бан/верификация данных)
    - Проверка портфолио моделей (модерация)
    - Управление заказами (возвраты)
    - Статистика:
        - Количество клиентов и моделей
        - Количество заказов
        - Средний чек 
        - Реферальные начисления
5. **Заказ услуги:**
    - Указывание адреса, времени
    - Мб доп услуги (фетиши у всех разные) - на усмотрение модели принимать или нет
    - Полный расчет стоимости заказа + % самой платформе
    - Статусы заказа:
        - В обработке
        - Подтвержден
        - В пути
        - Выполнен
        - Отменен
6. **Оплата:**
    - Банковские карты/кошельки
    - Автоматический перевод денежных средств модели после выполнения услуги
    - Возможен возврат по какой-то причине (указывается клиентом)
7. **Контроль качества:**
    - Оценки 1-5 звезд как моделей, так и клиентов
    - Текстовые отзывы моделям
8. **Поиск:**
    - Город
    - Пол модели
    - Дата и время
    - Стоимость
    - Категория «программы»
9. **Реферальная система**
   - У каждого пользователя (клиент/модель) при создании генерируется уникальный реферальный код (как в онлайн-казино, на букмекерских сайтах)
   - При регистрации новый юзер(клиент/модель) указывает данный код - «приглашение» 
      * любой юзер (клиент/модель) может пригласить до 7 юзеров
      * денежный бонус получает, как пригласивший (бонус зависит от того, в который раз приглашает), так и новый пользователь
10. Промокоды
    - Наименование 
    - % скидки 
    - Время действия (возможно всегда, как если промокод на др)
11. Программа лояльности
    - Бронзовый кролик (от 4 заказов): кэшбэк 2% 
    - Серебряный кролик (от 8 заказов): кэшбэк 3% 
    - Золотой кролик (от 12 заказов): кэшбэк 5% 
    - Платиновый кролик (от 20 заказов): кэшбэк 10%

## ER-диаграмма

```plantuml
entity "region" {
  * region_id : serial
  name : varchar(40)
}

entity "city" {
  * city_id : serial
  name : varchar(20)
  region : serial<<FK>>
}

"city" }o--|| "region" : "располагается"

entity "loyalty_level" {
  * level_id : serial
  name : varchar(30) // Бронзовый/Серебряный/Золотой/Платиновый кролик
  min_orders : integer
  cashback_percentage : integer
}

entity "auth" {
  * auth_id : serial
  email : varchar(320)
  phone : varchar(15)
  password_hash : varchar(255)
  created_at : timestamp
}

entity "admin" {
  * admin_id : serial
  auth_id : serial<<FK>>
  permissions : json // верификация/бан и тд
}

"auth" ||--|| "admin" : "вход"

entity "user" {
  * user_id : serial
  auth_id : serial<<FK>>
  birth_date : date
  gender : varchar(10) // фильтрация
  city_id : serial<<FK>> // фильтрация
  passport_series : varchar(4)
  passport_number : varchar(6)
  passport_issue_date : date
  passport_verified : boolean
  referral_code : uuid
  referral_user_id : serial<<FK>> // мб null
  referral_user_count : smallint
  is_banned : boolean
}

"auth" ||--|| "user" : "вход"
"user" }o--|| "city" : "находится"
"user" ||--o{ "user" : "от кого код-приглашение" 

' Бан пользователей (клиент/модель)
entity "ban" {
  * ban_id : serial
  admin_id : serial<<FK>>
  user_id : serial<<FK>>
  reason : varchar(500)
  created_at : timestamp
}

"admin" ||--o{ "ban" : "забанить"
"user" ||--o{ "ban" : "забаненный" // мб 0 банов

entity "client" {
  * client_id : serial
  user_id : serial<<FK>>
  name : varchar(50)
  loyalty_level_id : serial<<FK>>
}

"user" ||--|| "client" : "определяет"
"client" ||--o{ "loyalty_level" : "имеет"

entity "model" {
  * model_id : serial
  user_id : serial<<FK>>
  name : varchar(50)
}

"user" ||--|| "model" : "определяет"

entity "social_media" {
  * social_media_id : serial
  model_id : serial <<FK>>
  platform : varchar(50) // Instagram, Telegram и т. д.
  url : varchar(255)
}

"model" ||--o{ "social_media" : "публикуется"

entity "portfolio_data" {
  * portfolio_id : serial
  model_id : serial<<FK>>
  media_url : varchar(255) // фото/видео
  description : varchar(500)
  uploaded_at : timestamp
  is_verified : boolean
}

"model" ||--o{ "portfolio_data" : "имеет"

' Категория услуги (программа)
entity "category" {
  * category_id : serial
  name : varchar(70)
}

' Услуга
entity "service" {
  * service_id : serial
  category_id : serial<<FK>> // фильтрация
  description : varchar(255)
}

"service" }o--|| "category" : "принадлежит"

entity "model_service" {
  * model_service_id : serial
  model_id : serial<<FK>>
  service_id : serial<<FK>>
  price : decimal(9, 2)
}

"model" ||--o{ "model_service" : "предоставляет"
"service" ||--o{ "model_service" : "предоставляемая услуга"

' Забронировать/послать на одобрение услугу
entity "booking" {
  * booking_id : serial
  client_id : serial<<FK>>
  model_service_id : serial<<FK>>
  date_time : timestamp
  duration : timestamp
  address : json // частный дом/квартира/отель/заведение по типу ресторана и тд
  additional_service_id : serial<<FK>>
  status : varchar(20) // pending/rejected/approved/cancelled(by client)
  created_at : timestamp
}

"booking" }o--|| "model_service" : "выбор услуги"
"booking" }o--|| "client" : "забронировать"
"model" ||--o{ "booking" : "выбранная модель"

' Указание дополнительных услуг
entity "additional_service" {
  * additional_service_id : serial
  description : varchar(1000)
  offer_price : decimal(9, 2) // меняется как моделью, так и клиентом
  status : varchar(20) // pending/rejected/approved/cancelled(by client)/higherPrice(by model)
  updated_at : timestamp
}

"booking" ||--o{ "additional_service" : "имеет доп услуги"

entity "promocode" {
  * promocode_id : serial
  code : varchar(30)
  percentage : smallint
  start_time : timestamp
  finish_time : timestamp
  is_always : boolean // какой-то ежегодный праздник
}

' Заказ после одобрения услуги/ее брони
entity "order" {
  * order_id : serial
  booking_id : serial<<FK>>
  promocode_id : serial<<FK>>
  platform_fee : decimal(9, 2) // комиссия платформы
  total_cost : decimal(9, 2)
  status : varchar(20) // InProcess/Confirmed/InTransit/Completed/Cancelled
  created_at : timestamp
  security_code : varchar(6) // для безопасности - проверка при встрече
}

"order" ||--|| "booking" : "одобрение услуги"
"order" }o--o{ "promocode" : "применяется" // может быть и ноль промокодов

entity "payment_system_integration" {
  * payment_system_id : serial
  name : varchar(50) // Сбербанк/Тинькофф и тд
}

entity "external_transaction" {
  * external_transaction_id : serial
  payment_system_id : serial<<FK>>
  failure_msg : varchar(255) // платеж не прошел
}

"payment_system_integration" ||--o{ "external_transaction" : "обработка через"

entity "transaction" {
  * transaction_id : serial
  amount : decimal(9,2)
  type : varchar(20) // OrderPayment/OrderIncome/ClientDeposit/ModelPayout/RefundToClient/RefundFromModel/OrderCancellation/Referral/Сashback
  order_id : serial<<FK>> // мб null в случае ClientDeposit/ModelPayout/OrderCancellation/Referral
  external_transaction_id : serial<<FK>> // только в случае ClientDeposit/ModelPayout
  reason : varchar(255) // только в случае RefundToClient/RefundFromModel
  status : varchar(20) // failure/pending/success
  created_at : timestamp
  processed_at : timestamp
}

"order" }o--|| "transaction"  : "оплата заказа"
"transaction" }o--o{ "external_transaction"  : "внешняя транзакция" // мб 0

' Система отзывов - на их основе будет вычисляться рейтинг юзера
entity "review" {
  * review_id : serial
  order_id : serial<<FK>>
  from_user_id : serial<<FK>>
  to_user_id : serial<<FK>>
  rating : smallint // 1-5 целое число
  description : varchar(500)
  created_at : timestamp
}

"order" ||--o{ "review" : "оставить отзыв клиенту/модели по заказу"
"user" ||--o{ "review" : "от кого"
"user" ||--o{ "review" : "кому"

' Статистика - обновление ночью
entity "daily_statistics" {
  * stat_id : serial
  date : date
  total_clients : integer
  total_models : integer
  total_orders : integer
  completed_orders : integer // за день
  avg_order_cost : decimal(9,2) // avg(order.total_cost)
  total_referrals : integer
  referral_bonuses : decimal(9,2) // суммирование
}
```