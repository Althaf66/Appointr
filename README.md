## To start backend
```
docker-compose up --build
air
./stripe listen --forward-to localhost:8080/v1/webhook
```

## To start frontend
```
cd frontend
npm i
npm run dev
```

## To push to the github repo
```
git status
git add .
git commit --signoff -m "update"
git push origin main
```