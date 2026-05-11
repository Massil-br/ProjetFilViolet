from fastapi import FastAPI
from backend.api.routers.user_router import router as user_router
app = FastAPI()

app.include_router(user_router)



@app.get("/")
async def root():
    return {"message":"Welcome to pokerApi"}