from datetime import datetime
from decimal import Decimal
from sqlalchemy import Boolean, Enum, String,Numeric
from sqlalchemy import Text
from sqlalchemy import DateTime
from sqlalchemy import func


from sqlalchemy.orm import Mapped
from sqlalchemy.orm import mapped_column

from backend.api.database.db import Base


class UserRole(str, Enum):
    USER = "USER"
    MODERATOR = "MODERATOR"
    ADMIN = "ADMIN"



class User(Base):
    __tablename__ = "users"

    id: Mapped[int] = mapped_column(
        primary_key=True
    )

    firstname: Mapped[str] = mapped_column(
        String(32),
        nullable=False
    )
    lastname : Mapped[str] = mapped_column(
        String(32),
        nullable=False
    )
    nickname : Mapped[str] = mapped_column(
        String(50),
        nullable=False,
        unique=True
    )

    email: Mapped[str] = mapped_column(
        String(255),
        unique=True,
        nullable=False
    )

    hashed_password: Mapped[str] = mapped_column(
        Text,
        nullable=False
    )

    is_verified:Mapped[bool] = mapped_column(
        Boolean,
        default=False,
        nullable=False
    )

    status:Mapped[UserRole] = mapped_column(
        Enum(UserRole),
        default=UserRole.USER,
        nullable=False
    )

    is_online:Mapped[bool] = mapped_column(
        Boolean,
        default=False,
        nullable=False
    )

    money:Mapped[Decimal] = mapped_column(
        Numeric(10,2),
        nullable=False,
        default=500
    )

    created_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True),
        server_default=func.now()
    )