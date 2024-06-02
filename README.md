# Запуск сервиса для dev окружения

## Инструкции по установке

Для запуска с помощью Docker Compose, выполните следующие шаги:

1. Установите Docker и Docker Compose, если у вас их еще нет.

2. Склонируйте репозиторий Space Duck с GitHub:

   ```bash
   git clone https://github.com/kilievich-dmitriy-andreevich/violation-tracker-predictor-service.git
   ```

3. Перейдите в main каталог:

   ```bash
   cd violation-tracker-predictor-service
   ```

6. Запустите сервис с помощью Docker Compose:

   ```bash
   docker compose up -d
   ```

Теперь Space Duck должен быть доступен по адресу `http://localhost/docs`!


Если у вас ubuntu дистрибутив, то просто запустите заранее скомпиленный файл:

    sudo ./service -config_path ./configs/dev.yml