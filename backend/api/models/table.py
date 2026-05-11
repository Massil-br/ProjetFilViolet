from datetime import datetime
from decimal import Decimal
from sqlalchemy import Boolean,CheckConstraint, Enum, Integer, String,Numeric
from sqlalchemy import Text
from sqlalchemy import DateTime
from sqlalchemy import func
from sqlalchemy import ForeignKey


from sqlalchemy.orm import Mapped, relationship
from sqlalchemy.orm import mapped_column

from backend.api.database.db import Base


class Table(Base):
    __tablename__="tables"
    id:Mapped[int] = mapped_column(
        primary_key=True
    )
    saloon_id:Mapped[int] = mapped_column(
        ForeignKey("saloons.id"),
        nullable=False
    )
    slots_available:Mapped[int] = mapped_column(
        Integer
    )
    saloon = relationship("Saloon",back_populates="saloons")
    __table_args__ = (
        CheckConstraint("slots_available >= 2 AND slots_available <= 10", name="check_room_size")
    )