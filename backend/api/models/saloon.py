from datetime import datetime
from decimal import Decimal
from sqlalchemy import Boolean, Enum, Integer, String,Numeric
from sqlalchemy import Text
from sqlalchemy import DateTime
from sqlalchemy import func


from sqlalchemy.orm import Mapped
from sqlalchemy.orm import mapped_column

from backend.api.database.db import Base


class Saloon(Base):
    __tablename__="saloons"
    id:Mapped[int] = mapped_column(
        primary_key=True
    )
    tables_count: Mapped[int] = mapped_column(
        Integer,
        nullable=False,
        default=0
    )
    players_count: Mapped[int] = mapped_column(
        Integer,
        nullable=False,
        default=0
    )
    min_bet: Mapped[Decimal] = mapped_column(
        Numeric(10,2),
        default=10,
        nullable=False,
        min=10,
    )
    max_bet : Mapped[Decimal] = mapped_column(
        Numeric(10,2),
        nullable=False,
        default=50,
        min=10
    )
