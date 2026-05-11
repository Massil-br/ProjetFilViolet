from datetime import datetime
from decimal import Decimal
from sqlalchemy import Boolean, CheckConstraint, Enum, Integer, String,Numeric
from sqlalchemy import Text
from sqlalchemy import DateTime
from sqlalchemy import func
from sqlalchemy import ForeignKey


from sqlalchemy.orm import Mapped
from sqlalchemy.orm import mapped_column

from backend.api.database.db import Base

class Player(Base):
    __tablename__="players"
    user_id:Mapped[int] = mapped_column(
        ForeignKey("users.id"),
        primary_key=True,
        nullable=False

    )
    table_id:Mapped[int] = mapped_column(
        ForeignKey("tables.id"),
        primary_key=True,

    )
    initial_bet:Mapped[Decimal] = mapped_column(
        Numeric(10,2),
        nullable=False
    )
    money:Mapped[Decimal] = mapped_column(
        Numeric(10,2)
    )

    __table_args__ = (
        CheckConstraint("money >= 0",name="check_money_minimum")
    )