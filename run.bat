go build -o myapp.exe ./cmd/web/. || exit /b
myapp -dbname=bookings -dbuser=postgres -cache=false -production=false
myapp.exe