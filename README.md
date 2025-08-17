📌 Busha Recipient Service (Fiber + Golang)
This project provides an API endpoint to check if a Busha recipient exists, and if not, creates it using the Busha API.

🚀 Features
Checks existing recipients via Busha API

Creates a new recipient if not found

Environment variable configuration (no hardcoded credentials)

Fiber framework for fast HTTP handling

Go 1.22+ compatible (no deprecated imports)





📋 Requirements
Go 1.18 or higher (tested with Go 1.22)

Fiber v2

Busha Sandbox API account

Git installed

⚙️ Installation
1️⃣ Clone the repository


git clone https://github.com/your-username/busha-recipient-service.git
cd busha-recipient-service
2️⃣ Install dependencies


go mod tidy
3️⃣ Create .env file in the project root


BUSHA_API_KEY=your_busha_api_key
BUSHA_PROFILE_ID=your_busha_profile_id
BUSHA_API_VERSION=2025-07-11
Note:

Replace your_busha_api_key with the API key from Busha Sandbox.

Replace your_busha_profile_id with your profile ID.

Keep BUSHA_API_VERSION matching your API version.

4️⃣ Run the server


go run main.go
Server will start on:


http://localhost:3000
🛠 Usage

Create a Recipient

curl -X POST http://localhost:3000/recipient \
  -H "Content-Type: application/json" \
  -d '{
    "currency": "NGN",
    "country_code": "NG",
    "type": "ngn_bank",
    "bank_name": "Opay",
    "bank_code": "100004",
    "account_number": "9169277397",
    "account_name": "Ojieh Chukuwuyenum"
  }'


  
📜 Response Examples
If recipient exists


{
  "message": "Recipient already exists",
  "data": {
    "id": "12345",
    "account_number": "9169277397",
    "account_name": "Ojieh Chukuwuyenum",
    "bank_name": "Opay"
  }
}
If recipient is created


{
  "id": "67890",
  "account_number": "9169277397",
  "account_name": "Ojieh Chukuwuyenum",
  "bank_name": "Opay",
  "currency": "NGN",
  "country_code": "NG"
}
🧑‍💻 Development Notes
All secrets are stored in .env

Uses github.com/joho/godotenv to load environment variables

HTTP requests made with Go’s net/http package

No deprecated ioutil usage — replaced with io and os functions

