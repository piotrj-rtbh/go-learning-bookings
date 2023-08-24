go build -o bookings.exe ./cmd/web/.
bookings.exe -dbname=bookings -dbuser=postgres -dbpass=postgres -cache=false -production=false
