from datetime import datetime
from decimal import Decimal
from sqlalchemy import Boolean, Enum, String,Numeric
from sqlalchemy import Text
from sqlalchemy import DateTime
from sqlalchemy import func


from sqlalchemy.orm import Mapped
from sqlalchemy.orm import mapped_column

from api.database.db import Base
class Friend(Base):
    __tablename__ = "friends"
    id:Mapped[int] = mapped_column(
        primary_key=True
    )

