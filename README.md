# seaweed app

## backend

```bash
cd backend
go run cmd/app/main.go
```

needs `.env` with:
- `DB_ADDR` (neon postgres dsn)
- `FIREBASE_SERVICE_ACCOUNT_PATH` (path to json)
- `S3_SUBMISSIONS_BUCKET` (aws bucket name)
- `JUDGE0_BASE_URL`

## frontend

```bash
cd frontend
npm run dev
```

needs `.env.local` with:
- `NEXT_PUBLIC_API_URL` (usually http://localhost:8080)
- `NEXT_PUBLIC_FIREBASE_API_KEY`
- `NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN`
- `NEXT_PUBLIC_FIREBASE_PROJECT_ID`
- `NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET`
- `NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID`
- `NEXT_PUBLIC_FIREBASE_APP_ID`

## local testing

Set `DUMMY_AUTH=true` for backend and `NEXT_PUBLIC_AUTH_BYPASS=true` for frontend.

## kubernetes

Kubernetes manifests live in `seaweed-infra`.
