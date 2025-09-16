#!/bin/bash

if [ -f "cryptoserver" ]; then
    echo "Запуск скомпилированного crypto сервера..."
    ./cryptoserver
elif [ -f "cryptoserver.go" ]; then
    echo "Запуск Go crypto сервера..."
    go run cryptoserver.go
else
    echo "Не найден исполняемый файл crypto сервера"
    echo "Убедитесь что файл скомпилирован или существует cryptoserver.{py,js,go}"
    exit 1
fi