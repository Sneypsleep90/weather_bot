```markdown
# 🌦️ Простой Weather Bot (Go)

Telegram-бот, который моментально показывает текущую погоду в любом городе мира.

## ▶️ Как использовать
Просто отправьте боту **название города** (на русском или английском), и он пришлет актуальные погодные данные.

Пример:
```
Москва
```
или
```
London
```

## ⚙️ Технологии
- **Язык:** Go
- **Библиотека для Telegram:** [go-telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api/v5)
- **Погодный API:** OpenWeatherMap (или другой, который вы используете)

## 🚀 Установка и запуск

### 1. Клонирование репозитория
```bash
git clone https://github.com/ваш-логин/weather-bot-go.git
cd weather-bot-go
```

### 2. Настройка конфига
Создайте файл `.env` в корне проекта:
```ini
BOT_TOKEN=ваш_telegram_токен
OPENWEATHERAPI_KEY=ваш_ключ_openweather
```

### 3. Запуск
```bash
go run main.go
```

## 📊 Пример ответа бота
```
🌆 Город: Москва
🌡️ Температура: +15°C 
```

## 📦 Зависимости
Убедитесь, что у вас установлены:
- Go 1.20+
- Библиотеки (устанавливаются автоматически):
  ```bash
  go get github.com/go-telegram-bot-api/telegram-bot-api/v5
  go get github.com/joho/godotenv
  ```

