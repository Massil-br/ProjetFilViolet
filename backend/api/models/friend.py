from datetime import datetime
from sqlalchemy import ForeignKey
from decimal import Decimal
from sqlalchemy import Boolean, Enum, String,Numeric
from sqlalchemy import Text
from sqlalchemy import DateTime
from sqlalchemy import func


from sqlalchemy.orm import Mapped
from sqlalchemy.orm import mapped_column

from backend.api.database.db import Base
class Friend(Base):
    __tablename__ = "friends"
    user_id : Mapped[int] = mapped_column(
        ForeignKey("users.id"),primary_key=True
    )
    friend_id : Mapped[int] = mapped_column(ForeignKey("users.id"), primary_key=True)
    created_at: Mapped[datetime] = mapped_column(default=func.now)

