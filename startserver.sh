#!/bin/bash

# Настройка переменных окружения
export LC_ALL=C
export LANG=C

# Конфигурация SAMBA
SAMBA_IP="//192.168.53.5/snapshots"
HOST_IP=$(hostname -I | awk '{print $1}')

# Функция для логирования
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1"
}

log "Запуск сервера CBT Timelapses..."

# Монтирование SAMBA и мониторинг (в фоновом режиме)
monitor_samba_mount() {
    while true; do
        # Проверка, смонтирована ли папка
        if ! mountpoint -q /mnt/seefetch; then
            log "SAMBA (/mnt/seefetch) не смонтирована. Попытка монтирования..."
            sudo mount.cifs "$SAMBA_IP" /mnt/seefetch -o user=sambauser,pass=111111111
            
            if [ $? -eq 0 ]; then
                log "SAMBA успешно смонтирована."
            else
                log "ОШИБКА: Не удалось смонтировать SAMBA."
            fi
        fi
        # Проверка каждые 5 минут (300 секунд)
        sleep 300
    done
}

log "Запуск мониторинга SAMBA..."
monitor_samba_mount &
MOUNT_MONITOR_PID=$!

# Создание .env файла для frontend
log "Создание конфигурации frontend..."
touch /home/blunder/bin/cbt_timelapses_new/cbt_timelapses_frontend/.env
echo "VUE_APP_PATH_START=${HOST_IP}" > /home/blunder/bin/cbt_timelapses_new/cbt_timelapses_frontend/.env

# Функция для остановки всех процессов при завершении
cleanup() {
    log "Завершение работы сервера..."
    kill $MOUNT_MONITOR_PID 2>/dev/null
    kill $REDIS_PID 2>/dev/null
    kill $BACKEND_PID 2>/dev/null
    exit 0
}

# Установка обработчика сигналов
trap cleanup SIGINT SIGTERM

# Запуск Redis сервера
log "Запуск Redis сервера..."
redis-server --daemonize yes --bind 127.0.0.1 --port 6379
REDIS_PID=$(pgrep -f "redis-server")

if [ -n "$REDIS_PID" ]; then
    log "Redis сервер запущен (PID: $REDIS_PID, порт: 6379)"
else
    log "ОШИБКА: Не удалось запустить Redis сервер"
    exit 1
fi

# Пауза для инициализации Redis
sleep 1

# Запуск Backend (Go сервер на порту 5000)
log "Запуск Backend сервера..."
cd /home/blunder/bin/cbt_timelapses_new/cbt_timelapses_backend
./main &
BACKEND_PID=$!

if kill -0 $BACKEND_PID 2>/dev/null; then
    log "Backend сервер запущен (PID: $BACKEND_PID, порт: 5000)"
else
    log "ОШИБКА: Не удалось запустить Backend сервер"
    exit 1
fi

# Сборка Frontend для продакшена (если нужно)
log "Проверка сборки Frontend..."
cd /home/blunder/bin/cbt_timelapses_new/cbt_timelapses_frontend

if [ ! -d "dist" ] || [ ! -f "dist/index.html" ]; then
    log "Собираем Frontend..."
    npm run build
    if [ $? -ne 0 ]; then
        log "ОШИБКА: Не удалось собрать Frontend"
        kill $BACKEND_PID 2>/dev/null
        exit 1
    fi
else
    log "Frontend уже собран"
fi

log "Все сервера успешно запущены!"
log "Redis сервер доступен на: 127.0.0.1:6379"
log "Приложение доступно на: http://${HOST_IP}:5000"
log "Backend API и Frontend интегрированы на одном порту!"
log "Для остановки нажмите Ctrl+C"

# Ожидание завершения процессов
wait $BACKEND_PID
