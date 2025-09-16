#!/bin/bash

if [ -f "cryptoserver.go" ]; then
    echo "Компиляция Go crypto сервера..."
    go build -o cryptoserver cryptoserver.go
else
    echo "Не найден файл cryptoserver для компиляции"
    echo "Поддерживаемые файлы: cryptoserver.{cpp,go,py,js,java}"
    exit 1
fi

echo "Компиляция cryptoserver завершена"