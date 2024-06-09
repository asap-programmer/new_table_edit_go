# practGo
Таблица с редактированием в MySQL на Go (Golang) c триггером и хранимой процедурой

Запустить сервер через run_select.bat, перейти по ссылке: http://localhost:8181/
На экране будет примерно следующее:

![Screenshot_1](https://github.com/asap-programmer/new_table_edit_go/assets/123025209/a154da82-fbaf-4cfb-886f-a51ce14a714d)

Сайт на Go (Golang) гарантированно работает при установленных программах и компонентах под Windows 11:
1) Go версии: go1.20.2 windows/amd64 в каталог: c:\Go\bin (из файла: go1.20.2.windows-amd64.msi);
2) Git версии: 2.40.0.windows.1 (из файла: Git-2.40.0-64-bit.exe);
3) Go SQL Driver с помощью команды: go get -u github.com/go-sql-driver/mysql .
